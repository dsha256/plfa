package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	PragmaticFeedWsURL  = "PRAGMATIC_FEED_WS_URL"
	CasinoID            = "CASINO_ID"
	TableIDs            = "TABLE_IDS"
	CurrencyIDs         = "CURRENCY_IDS"
	ServerPort          = "SERVER_PORT"
	PusherChannelID     = "PUSHER_CHANNEL_ID"
	PusherPeriodMinutes = "PUSHER_PERIOD_MINUTES"
	PusherAppID         = "PUSHER_APP_ID"
	PusherKey           = "PUSHER_KEY"
	PusherSecret        = "PUSHER_SECRET"
	PusherCluster       = "PUSHER_CLUSTER"
)

type ENV struct {
	PragmaticFeedWsURL  string
	CasinoID            string
	TableIDs            string
	CurrencyIDs         string
	ServerPort          string
	PusherChannelID     string
	PusherPeriodMinutes string
	PusherAppID         string
	PusherKey           string
	PusherSecret        string
	PusherCluster       string
}

func (env *ENV) Load() *ENV {
	env.PragmaticFeedWsURL = os.Getenv(PragmaticFeedWsURL)
	env.CasinoID = os.Getenv(CasinoID)
	env.TableIDs = os.Getenv(TableIDs)
	env.CurrencyIDs = os.Getenv(CurrencyIDs)
	env.ServerPort = os.Getenv(ServerPort)
	env.PusherChannelID = os.Getenv(PusherChannelID)
	env.PusherPeriodMinutes = os.Getenv(PusherPeriodMinutes)
	env.PusherAppID = os.Getenv(PusherAppID)
	env.PusherKey = os.Getenv(PusherKey)
	env.PusherSecret = os.Getenv(PusherSecret)
	env.PusherCluster = os.Getenv(PusherCluster)

	return env
}

// TODO: Reflection can make your life easier. Give this present to yourself this new year :)

func (env *ENV) GetPragmaticFeedWsURL() string {
	panicOnEmptyEnvVar(PragmaticFeedWsURL, env.PragmaticFeedWsURL)
	return strCleanUp(env.PragmaticFeedWsURL)
}

func (env *ENV) GetTableIDs() []string {
	res := strings.Split(strCleanUp(env.TableIDs), ",")
	panicOnEmptyEnvVar(TableIDs, res)
	return res
}

func (env *ENV) GetCurrencyIDs() []string {
	res := strings.Split(strCleanUp(env.CurrencyIDs), ",")
	panicOnEmptyEnvVar(CurrencyIDs, res)
	return res
}

func (env *ENV) GetCasinoID() string {
	panicOnEmptyEnvVar(CasinoID, env.CasinoID)
	return strCleanUp(env.CasinoID)
}

func (env *ENV) GetServerPort() string {
	panicOnEmptyEnvVar(ServerPort, env.ServerPort)
	return strCleanUp(env.ServerPort)
}

func (env *ENV) GetPusherChannelID() string {
	panicOnEmptyEnvVar(PusherChannelID, env.PusherChannelID)
	return strCleanUp(env.PusherChannelID)
}

func (env *ENV) GetPusherPeriodMinutes() int {
	cleanedUpMinutes := strCleanUp(env.PusherPeriodMinutes)
	minutes, err := strconv.Atoi(cleanedUpMinutes)
	if err != nil {
		log.Fatalf("Error convertin %s to int", PusherPeriodMinutes)
	}
	panicOnEmptyEnvVar(PusherPeriodMinutes, minutes)
	return minutes
}

func (env *ENV) GetPusherAppID() string {
	panicOnEmptyEnvVar(PusherAppID, env.PusherAppID)
	return strCleanUp(env.PusherAppID)
}

func (env *ENV) GetPusherKey() string {
	panicOnEmptyEnvVar(PusherKey, env.PusherKey)
	return strCleanUp(env.PusherKey)
}

func (env *ENV) GetPusherSecret() string {
	panicOnEmptyEnvVar(PusherSecret, env.PusherSecret)
	return strCleanUp(env.PusherSecret)
}

func (env *ENV) GetPusherCluster() string {
	panicOnEmptyEnvVar(PusherCluster, env.PusherCluster)
	return strCleanUp(env.PusherCluster)
}
