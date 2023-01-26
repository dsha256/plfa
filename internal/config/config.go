package config

import (
	"github.com/spf13/viper"
)

const (
	configName = "app"
	configType = "env"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	PragmaticFeedWSURL  string `mapstructure:"PRAGMATIC_FEED_WS_URL"`
	CasinoID            string `mapstructure:"CASINO_ID"`
	TableIDs            string `mapstructure:"TABLE_IDS"`
	CurrencyIDs         string `mapstructure:"CURRENCY_IDS"`
	ServerPort          string `mapstructure:"SERVER_PORT"`
	PusherChannelID     string `mapstructure:"PUSHER_CHANNEL_ID"`
	PusherPeriodMinutes int    `mapstructure:"PUSHER_PERIOD_MINUTES"`
	PusherAppID         string `mapstructure:"PUSHER_APP_ID"`
	PusherKey           string `mapstructure:"PUSHER_KEY"`
	PusherSecret        string `mapstructure:"PUSHER_SECRET"`
	PusherCluster       string `mapstructure:"PUSHER_CLUSTER"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
