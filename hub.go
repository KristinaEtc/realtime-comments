package main

import (
	"time"
)

type hub struct {
	// Registered clients
	clients map[*client]bool

	// Inbound messages
	broadcast chan string

	// Register requests
	register chan *client

	// Unregister requests
	unregister chan *client

	content string
}

var h = hub{
	broadcast:  make(chan string),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
	content:    "",
}

func sendPerioticData(c *client) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			log.Debug("gg")
			data = append([]byte(time.Now().String()), data...)
			c.send <- []byte(data)
			//	err := c.write(websocket.PingMessage, data)
			//	if err != nil {
			//		log.Errorf("Error write to client: %s", err.Error())
			//	}
		}
	}
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			log.Debugf("send <- []byte(h.content): [%v]", h.content)
			c.send <- []byte(h.content)
			go sendPerioticData(c)
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break

		case m := <-h.broadcast:
			h.content = m
			h.broadcastMessage()
			break
		}
	}
}

func (h *hub) broadcastMessage() {
	for c := range h.clients {
		select {
		case c.send <- []byte(h.content):
			log.Debugf("content=[%s]", h.content)
			break

		// We can't reach the client
		default:
			log.Debug("closing broadcast")
			close(c.send)
			delete(h.clients, c)
		}
	}
}
