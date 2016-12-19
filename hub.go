package main

import "log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// all history since server was run
	savedMessages [][]byte
}

func newHub() *Hub {
	return &Hub{
		broadcast:     make(chan []byte),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
		savedMessages: make([][]byte, 0),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			// give history back
			// TODO: only give history back after verify
			for _, msg := range h.savedMessages {
				client.send <- msg
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {

				msg := Message{disconnect, client.username, "has disconnected"}.encode()

				go func() {
					h.broadcast <- msg
				}()

				log.Println(client.username, "disconnected")

				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:

			// save to history
			h.savedMessages = append(h.savedMessages, message)

			msg := newMessage(message)
			log.Printf("%v: %v", msg.User, msg.Msg)

			// TODO: check if sent conn message and only broadcast if yes

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
