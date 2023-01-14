package main

import (
	"encoding/json"
	"fmt"
	"github.com/dsha256/plfa/internal/config"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/dsha256/plfa/internal/ws"
	"github.com/dsha256/plfa/pkg/dto"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	eventTmpl = "{\"type\":\"subscribe\",\"key\":\"%s\",\"casinoId\":\"%s\",\"currency\":\"%s\"}"
)

var errChan = make(chan dto.Error)
var dataChan = make(chan []byte)

func main() {
	bootstrap()
}

func bootstrap() {
	env := config.ENV{}
	env.Load()

	sigChan := make(chan os.Signal, 1)

	// ref: https://pkg.go.dev/os/signal
	signal.Notify(sigChan,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	var conns []*ws.Client

	// WebSocket Clients for each msh since the Pragma WS API allows us to write only one msg in one client.
	eventMsgs := genEventMsgs(env.GetTableIDs(), env.GetCurrencyIDs(), env.GetCasinoID())
	for _, msg := range eventMsgs {
		wsConn, _, _ := websocket.DefaultDialer.Dial(env.GetPragmaticFeedWsURL(), nil)
		wsClient := ws.NewClient(wsConn, errChan, dataChan)
		conns = append(conns, wsClient)
		go wsClient.WriteMsg([]byte(msg))
		//go wsClient.ReadMsgs()
	}

	aggregatedRepo := repository.NewAggregator()

	// React on errors and system signals.
	var pt dto.PragmaticTable
	for {
		select {
		case data := <-dataChan:
			err := json.Unmarshal(data, &pt)
			if err != nil {
				panic(err)
			}
			aggregatedRepo.AddTable(pt)
		case err := <-errChan:
			log.Println(err)
		case sig := <-sigChan:
			switch sig {
			case os.Interrupt:
				log.Println("Interrupted. Shutting down gracefully...")
				wsConsShutdown(conns...)
				select {
				case <-time.After(5 * time.Second):
					log.Println("Gracefully shut down.")
				}
				closeChans()
				os.Exit(0)
			case syscall.SIGTERM:
				log.Println("SIGTERM. Shutting down gracefully...")
				wsConsShutdown(conns...)
				select {
				case <-time.After(5 * time.Second):
					log.Println("Gracefully shut down.")
				}
				closeChans()
				os.Exit(0)
			case syscall.SIGQUIT:
				log.Println("SIGQUIT. Shutting down gracefully...")
				wsConsShutdown(conns...)
				select {
				case <-time.After(5 * time.Second):
					log.Println("Gracefully shut down.")
				}
				closeChans()
				os.Exit(0)
			}
		}
	}
}

func wsConsShutdown(conns ...*ws.Client) {
	for _, conn := range conns {
		conn.ShutDown()
	}
}

func genEventMsgs(tbIDs, curIDs []string, casinoID string) []string {
	var eventMsgs []string
	for _, tbID := range tbIDs {
		for _, curID := range curIDs {
			eventMsgs = append(eventMsgs, fmt.Sprintf(eventTmpl, tbID, casinoID, curID))
		}
	}
	return eventMsgs
}

func closeChans() {
	close(errChan)
	close(dataChan)
}
