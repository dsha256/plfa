package pusher

import (
	"fmt"
	"github.com/dsha256/plfa/internal/jsonlog"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/pusher/pusher-http-go/v5"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// channelName is a Pusher's channel name.
	channelName = "update-pragmatic-live-feed-tables"
	// eventName is a Pusher's event name
	eventName = "LiveFeedTables"
)

type Client struct {
	client                *pusher.Client
	channelID             string
	pusherPeriodInMinutes int
	repo                  repository.AggregateRepository
	logger                *jsonlog.Logger
}

type Service struct {
	ChannelID            string
	PushingPeriodMinutes int
}

type Secrets struct {
	AppID   string
	Key     string
	Secret  string
	Cluster string
	Secure  bool
}

type Configs struct {
	Service Service
	Secrets Secrets
	Repo    repository.AggregateRepository
	Logger  *jsonlog.Logger
}

func NewClient(cfg *Configs) *Client {

	c := pusher.Client{
		AppID:   cfg.Secrets.AppID,
		Key:     cfg.Secrets.Key,
		Secret:  cfg.Secrets.Secret,
		Cluster: cfg.Secrets.Cluster,
		Secure:  cfg.Secrets.Secure,
	}

	return &Client{
		client:                &c,
		channelID:             cfg.Service.ChannelID,
		pusherPeriodInMinutes: cfg.Service.PushingPeriodMinutes,
		repo:                  cfg.Repo,
		logger:                cfg.Logger,
	}
}

func (c *Client) StartPushing() error {
	ticker := time.NewTicker(time.Duration(c.pusherPeriodInMinutes) * time.Minute)
	defer ticker.Stop()

	c.logger.PrintInfo("starting pusher service", map[string]string{
		"pushing_period": fmt.Sprintf("%d minute(s)", c.pusherPeriodInMinutes),
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			pragmaticTables, err := c.repo.ListTables()
			if err != nil {
				c.logger.PrintError(err, nil)
				return err
			}
			err = c.client.Trigger(channelName, eventName, pragmaticTables)
			if err != nil {
				c.logger.PrintError(err, nil)
				return err
			}
			c.logger.PrintInfo("successfully pushed table's updated data", nil)

		case sig := <-quit:
			c.logger.PrintInfo("stopped pusher service", map[string]string{"signal": sig.String()})
			return nil
		}
	}
}
