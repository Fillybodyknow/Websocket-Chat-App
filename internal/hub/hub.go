package hub

import (
	"log"
)

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

type Message struct {
	Data      []byte
	RoomID    string
	SenderRef *Client
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.RoomID] == nil {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.Rooms[client.RoomID][client] = true
			log.Println("🔵 Client joined room:", client.RoomID)

		case client := <-h.Unregister:
			if clients, ok := h.Rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					log.Println("🔴 Client left room:", client.RoomID)
				}
			}

		case msg := <-h.Broadcast:
			if clients, ok := h.Rooms[msg.RoomID]; ok {
				for client := range clients {
					if client != msg.SenderRef {
						select {
						case client.Send <- msg.Data:
						default:
							close(client.Send)
							delete(clients, client)
						}
					}
				}
			}

		}
	}
}
