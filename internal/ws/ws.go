package ws

import (
	"errors"
	"github.com/dsha256/plfa/internal/jsonlog"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/dsha256/plfa/pkg/dto"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serverConnCloseTimeout = 1 * time.Second

var (
	interrupt = make(chan os.Signal, 1)
	done      = make(chan struct{})

	// clients is started clients quantity.
	clients = 0
)

type Client struct {
	repo   repository.AggregateRepository
	logger *jsonlog.Logger
}

func NewClient(repo repository.AggregateRepository, logger *jsonlog.Logger) *Client {
	return &Client{repo: repo, logger: logger}
}

// RunAndListenClient creates and runs WS client based on the url passed in as a parameter. Then it writes the message
// passed as a second parameter and writes incoming payload to the repository.
func (c *Client) RunAndListenClient(url string, msg string) {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		c.logger.PrintFatal(err, nil)
	} else {
		clients++
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			c.logger.PrintError(err, nil)
		}
	}(conn)

	go func() {
		defer func() {
			clients--
			// Avoiding closing of closed channel.
			if clients == 0 {
				close(done)
				return
			}
			done <- struct{}{}
		}()
		for {
			_, payload, err := conn.ReadMessage()
			if err != nil {
				// An appropriate client is closed, so we can just stop the reader.
				if errors.As(err, &websocket.ErrCloseSent) {
					return
				}
				c.logger.PrintError(err, nil)
				return
			}
			pragmaticTb, err := dto.Bytes2PT(payload)
			if err != nil {
				c.logger.PrintFatal(err, nil)
			}
			c.repo.AddTable(pragmaticTb)
		}
	}()

	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		//log.Println("write 59:", err)
		return
	}
	c.logger.PrintInfo("listening new WS", map[string]string{
		"local_address":  conn.LocalAddr().String(),
		"remote_address": conn.RemoteAddr().String(),
	})

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-done:
			c.logger.PrintInfo("stopped WS clients", nil)
			return
		case sig := <-interrupt:
			c.logger.PrintInfo("shutting down WS clients", map[string]string{"signal": sig.String()})

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				c.logger.PrintError(err, nil)
				return
			}
			select {
			case <-done:
			case <-time.After(serverConnCloseTimeout):
				c.logger.PrintError(errors.New("timeout during closing WS connection"), nil)
				return
			}
			return
		}
	}
}
