package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type client struct {
	ws    *websocket.Conn
	send  chan []byte // Channel storing outcoming messages
	close chan bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		log.Errorf("[%s]: Method not allowed ", r.RemoteAddr)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Web-scoket serve: %s", err.Error())
		log.Errorf("[%s]: Web-scoket serve: %s ", r.RemoteAddr, err.Error())
		return
	}

	c := &client{
		send:  make(chan []byte, maxMessageSize),
		ws:    ws,
		close: make(chan bool),
	}

	h.register <- c

	go c.writePump()
	c.readPump()
}

func (c *client) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(int64(maxMessageSize))
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Errorf("[%s] could not read from web-socket: [%s]", c.ws.RemoteAddr(), err.Error())
			break
		}

		log.Debugf("[%s] got message [%s]", c.ws.RemoteAddr(), message)
		if strings.TrimRight(string(message), "\n") == globalOpt.MonitoringMessage {
			log.Infof("[%s] set monitoring client", c.ws.RemoteAddr())
			h.setMonitoringClient(c)
		}

		msgTime := time.Now().Format("[2006-01-02/15:04:05] ")
		message = append([]byte(msgTime), message...)

		if globalOpt.Broadcast {
			h.broadcast <- string(message)
		} else {
			c.send <- message
			if h.monitoringClient != nil {
				h.monitoringClient.send <- message
			}
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Errorf("[%s] could not get message. Closing web-scokets", c.ws.RemoteAddr())
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Errorf("[%s] could not write message [%s]", c.ws.RemoteAddr(), err.Error())
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("[%s] could not write periodic data [%s]", c.ws.RemoteAddr(), err.Error())
				return
			}
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}
