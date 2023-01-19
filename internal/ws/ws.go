package ws

import (
	"errors"
	"github.com/dsha256/plfa/internal/repository"
	"github.com/dsha256/plfa/pkg/dto"
	"github.com/gorilla/websocket"
	"log"
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

func init() {
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
}

type Client struct {
	repo repository.AggregateRepository
}

func NewClient(repo repository.AggregateRepository) *Client {
	return &Client{repo: repo}
}

// RunAndListenClient creates and runs WS client based on the url passed in as a parameter. Then it writes the message
// passed as a second parameter and writes incoming payload to the repository.
func (client *Client) RunAndListenClient(url string, msg string) {

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial 26:", err)
	} else {
		clients++
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println("Error closing WS connection.")
		}
	}(c)

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
			_, payload, err := c.ReadMessage()
			if err != nil {
				if errors.As(err, &websocket.ErrCloseSent) {
					return
				}
				log.Println("read: 52222", err)
				return
			}
			pragmaticTb, err := dto.Bytes2PT(payload)
			if err != nil {
				log.Println(err)
			}
			client.repo.AddTable(pragmaticTb)
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write 59:", err)
		return
	}

	for {
		select {
		case <-done:
			return
		case s := <-interrupt:
			log.Println(s.String())

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(serverConnCloseTimeout):
				log.Println("Timeout during closing WS connection.")
				return
			}
			return
		}
	}
}
