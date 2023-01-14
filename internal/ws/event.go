package ws

import "encoding/json"

type Event struct {
	// Type is the message type to sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

const (
	// EventSendMessage is the event name for new chat messages sent
	EventSendMessage = "send_message"
	// EventNewMessage is a response to send_message
	EventNewMessage = "new_message"
	// EventChangeRoom is event when switching rooms
	EventChangeRoom = "change_room"
)
