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
	ticker := time.NewTicker(time.Second * time.Duration(globalOpt.WriteTestDataTimeout))
	for {
		select {
		case _ = <-c.close:
			log.Debugf("[%s] closing client channel", c.ws.RemoteAddr())
			return
		case t := <-ticker.C:
			timeMesg := t.Format("[2006-01-02/15:04:05] ")
			log.Debugf("[%s] sending test data %s", c.ws.RemoteAddr(), timeMesg)
			msg := append([]byte(timeMesg), data...)
			c.send <- []byte(msg)
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
			log.Infof("[%s] register new client", c.ws.RemoteAddr())
			c.send <- []byte(h.content)
			go sendPerioticData(c)
			break

		case c := <-h.unregister:
			log.Infof("[%s] unregister new client", c.ws.RemoteAddr())
			_, ok := h.clients[c]
			if ok {
				c.close <- true
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
			break

		// We can't reach the client
		default:
			log.Debug("closing broadcast")
			close(c.send)
			delete(h.clients, c)
		}
	}
}
