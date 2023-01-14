package ws

import "github.com/gorilla/websocket"

type Manager struct {
	clients ClientList
	conn    *websocket.Conn
}
