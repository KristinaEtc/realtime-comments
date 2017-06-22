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
	maxMessageSize = 2056
)

var connectedCount int
var monitoringClient *client

type client struct {
	ws     *websocket.Conn
	sendCh chan []byte
	close  chan bool
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *client) read() {
	log.Debug("read start")

	defer func() {
		log.Debug("read exiting")
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Errorf("[%s] could not read from web-socket: [%s]", c.ws.RemoteAddr(), err.Error())
			if monitoringClient == c {
				monitoringClient = nil
			}
			return
		}

		log.Debugf("[%s] got message [%s]", c.ws.RemoteAddr(), message)
		if strings.TrimRight(string(message), "\n") == globalOpt.MonitoringMessage {
			log.Infof("[%s] set monitoring client", c.ws.RemoteAddr())
			monitoringClient = c
		}

		msgTime := time.Now().Format("[2006-01-02/15:04:05] ")
		message = append(message, data...)
		message = append([]byte(msgTime), message...)
		c.sendCh <- message

		if monitoringClient != nil && c != monitoringClient {
			monitoringClient.sendCh <- message
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}

func (c *client) process() {
	log.Debug("writePump start")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Debug("processWrite close")
		ticker.Stop()
		connectedCount--
		if monitoringClient == c {
			monitoringClient = nil
		}
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendCh:
			if !ok {
				log.Errorf("[%s] could not get message. Closing web-scokets", c.ws.RemoteAddr())
				c.write(websocket.CloseMessage, []byte{})
				if monitoringClient == c {
					monitoringClient = nil
				}
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Errorf("[%s] could not write message [%s]", c.ws.RemoteAddr(), err.Error())
				if monitoringClient == c {
					monitoringClient = nil
				}
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("[%s] could not write periodic data [%s]", c.ws.RemoteAddr(), err.Error())
				if monitoringClient == c {
					monitoringClient = nil
				}
				return
			}
		case _ = <-c.close:
			log.Debugf("[%s] Closing web-scokets", c.ws.RemoteAddr())
			if monitoringClient == c {
				monitoringClient = nil
			}
			return
		}
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Debug("serveWs")

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
		sendCh: make(chan []byte, maxMessageSize),
		ws:     ws,
	}

	go c.process()
	c.read()
}
