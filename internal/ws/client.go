package ws

import (
	"github.com/dsha256/plfa/pkg/dto"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const serviceName = "WS"
const writeInterval = 1 * time.Second

type ClientList map[*Client]bool

type Client struct {
	Conn     *websocket.Conn
	ErrChan  chan dto.Error
	DataChan chan []byte
}

func NewClient(conn *websocket.Conn, errChan chan dto.Error, dataChan chan []byte) *Client {
	return &Client{Conn: conn, ErrChan: errChan, DataChan: dataChan}
}

func (c *Client) WriteMsg(msg []byte) {

	ticker := time.NewTicker(writeInterval)
	defer ticker.Stop()

	go c.ReadMsgs()
	for {
		select {
		case <-ticker.C:
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				e := dto.Error{ServiceName: serviceName, Err: err}
				c.ErrChan <- e
			}
		}
	}
}

func (c *Client) ReadMsgs() {
	for {
		_, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				e := dto.Error{ServiceName: serviceName, Err: err}
				c.ErrChan <- e
				return
			}
			e := dto.Error{ServiceName: serviceName, Err: err}
			c.ErrChan <- e
			return
		}
		c.DataChan <- payload
	}
}

func (c *Client) ShutDown() {
	log.Printf("Starting graceful shutdown: %s\n", c.Conn.LocalAddr().String())
	err := c.Conn.WriteMessage(websocket.CloseMessage, nil)
	if err != nil {
		e := dto.Error{ServiceName: serviceName, Err: err}
		c.ErrChan <- e
		return
	}
	err = c.Conn.Close()
	if err != nil {
		e := dto.Error{ServiceName: serviceName, Err: err}
		c.ErrChan <- e
		return
	}
	log.Printf("Finished gracefull shutdown: %s\n", c.Conn.LocalAddr().String())
}
