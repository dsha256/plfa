package main

import (
	"fmt"
	"github.com/dsha256/plfa/internal/config"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/dsha256/plfa/internal/ws"
	"sync"
)

const (
	// wsMsgTmpl is a template for creating message to write into the Web Socket connection.
	wsMsgTmpl = "{\"type\":\"subscribe\",\"key\":\"%s\",\"casinoId\":\"%s\",\"currency\":\"%s\"}"
)

func main() {
	bootstrap()
}

func bootstrap() {
	env := config.ENV{}
	env.Load()

	wg := sync.WaitGroup{}

	aggregatedRepo := repository.NewAggregator()

	msgs := genMsgs(env.GetTableIDs(), env.GetCurrencyIDs(), env.GetCasinoID())
	wsClient := ws.NewClient(aggregatedRepo)
	for _, msg := range msgs {
		wg.Add(1)
		go func(msg string) {
			defer wg.Done()
			wsClient.RunAndListenClient(env.GetPragmaticFeedWsURL(), msg)
		}(msg)
	}

	wg.Wait()
}

// genMsgs generates messages based on wsMsgTmpl format.
func genMsgs(tbIDs, curIDs []string, casinoID string) []string {
	var eventMsgs []string
	for _, tbID := range tbIDs {
		for _, curID := range curIDs {
			eventMsgs = append(eventMsgs, fmt.Sprintf(wsMsgTmpl, tbID, casinoID, curID))
		}
	}
	return eventMsgs
}
