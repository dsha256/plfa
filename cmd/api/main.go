package main

import (
	"fmt"
	"github.com/dsha256/plfa/internal/config"
	"github.com/dsha256/plfa/internal/jsonlog"
	"github.com/dsha256/plfa/internal/pusher"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/dsha256/plfa/internal/server"
	"github.com/dsha256/plfa/internal/ws"
	"os"
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

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	wg := sync.WaitGroup{}

	aggregatedRepo := repository.NewAggregator()

	// Web Socket Client.
	msgs := genMsgs(env.GetTableIDs(), env.GetCurrencyIDs(), env.GetCasinoID())
	// TODO: refactor for using custom logger
	wsClient := ws.NewClient(aggregatedRepo)
	for _, msg := range msgs {
		wg.Add(1)
		go func(msg string) {
			defer wg.Done()
			wsClient.RunAndListenClient(env.GetPragmaticFeedWsURL(), msg)
		}(msg)
	}

	// Server.
	newServer := server.NewServer(logger, aggregatedRepo)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := newServer.Serve(env.GetServerPort())
		if err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	// Pusher Client.
	pusherCfg := pusher.Configs{
		Service: pusher.Service{
			ChannelID:            env.GetPusherChannelID(),
			PushingPeriodMinutes: env.GetPusherPeriodMinutes(),
		},
		Secrets: pusher.Secrets{
			AppID:   env.GetPusherAppID(),
			Key:     env.GetPusherKey(),
			Secret:  env.GetPusherSecret(),
			Cluster: env.GetPusherCluster(),
			Secure:  true,
		},
		Repo:   aggregatedRepo,
		Logger: logger,
	}
	pusherClient := pusher.NewClient(&pusherCfg)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := pusherClient.StartPushing()
		if err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	// Wait for all the services.
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
