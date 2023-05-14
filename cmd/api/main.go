package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dsha256/plfa/internal/config"
	"github.com/dsha256/plfa/internal/jsonlog"
	"github.com/dsha256/plfa/internal/pusher"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/dsha256/plfa/internal/server"
	"github.com/dsha256/plfa/internal/ws"
)

const (
	// wsMsgTmpl is a template for creating message to write into the Web Socket connection.
	wsMsgTmpl = "{\"type\":\"subscribe\",\"key\":\"%s\",\"casinoId\":\"%s\",\"currency\":\"%s\"}"

	configFilePath = "."
)

// @title Pragmatic Live Feed Aggregator API Documentation
// @version 1.0.0
// @host localhost:8080
// @BasePath /v1
func main() {
	bootstrap()
}

func bootstrap() {

	// TODO: refactor env-related logic for delivery workflow simplicity.
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("can not load config")
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	aggregatedRepo := repository.NewAggregator()

	wg := sync.WaitGroup{}

	msgs := genMsgs(str2slice(cfg.TableIDs, ","), str2slice(cfg.CurrencyIDs, ","), cfg.CasinoID)
	wsClient := ws.NewClient(aggregatedRepo, logger)
	for _, msg := range msgs {
		wg.Add(1)
		go func(msg string) {
			defer wg.Done()
			wsClient.RunAndListenClient(cfg.PragmaticFeedWSURL, msg)
		}(msg)
	}

	newServer := server.NewServer(logger, aggregatedRepo)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := newServer.Serve(cfg.ServerPort)
		if err != nil {
			logger.PrintFatal(err, nil)
		}
	}()

	pusherCfg := pusher.Configs{
		Service: pusher.Service{
			ChannelID:            cfg.PusherChannelID,
			PushingPeriodMinutes: cfg.PusherPeriodMinutes,
		},
		Secrets: pusher.Secrets{
			AppID:   cfg.PusherAppID,
			Key:     cfg.PusherKey,
			Secret:  cfg.PusherSecret,
			Cluster: cfg.PusherCluster,
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

func str2slice(str, sep string) []string {
	if strings.HasSuffix(str, sep) {
		str = str[:len(str)-2]
	}
	return strings.Split(str, sep)
}
