package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/KristinaEtc/realtime-comments/database"
	"github.com/gorilla/websocket"
	"github.com/ventu-io/slf"
)

//TODO: move to configuration file
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
	log    slf.Logger
	db     database.Database
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *client) processCommentData(comment []byte) {
	c.log.Debugf("got message [%s]", comment)
	if strings.TrimRight(string(comment), "\n") == globalOpt.ServerConfig.MonitoringMessage {
		c.log.Info("set monitoring client")
		monitoringClient = c
	}

	currTime := time.Now()

	msgTime := currTime.Format("[2006-01-02/15:04:05] ")
	comment = append(comment, data...)
	comment = append([]byte(msgTime), comment...)
	go c.db.InsertData(comment, currTime)

	lastComment, err := c.db.GetData()
	if err != nil {
		c.log.Errorf("get data from [%s]: [%s]", globalOpt.DatabaseConfig.Type, err.Error())
	}
	c.sendCh <- lastComment
	if monitoringClient != nil && c != monitoringClient {
		monitoringClient.sendCh <- comment
	}
}

func (c *client) read() {
	c.log.Debug("read start")

	defer func() {
		c.log.Debug("read exiting")
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			c.log.Errorf("could not read from web-socket: [%s]", err.Error())
			if monitoringClient == c {
				monitoringClient = nil
			}
			return
		}
		c.processCommentData(message)
	}
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}

func (c *client) process() {
	c.log.Debug("writePump start")
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		c.log.Debug("processWrite close")
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
				c.log.Error("could not get message. Closing web-scokets")
				c.write(websocket.CloseMessage, []byte{})
				if monitoringClient == c {
					monitoringClient = nil
				}
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				c.log.Errorf("could not write message [%s]", err.Error())
				if monitoringClient == c {
					monitoringClient = nil
				}
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				c.log.Errorf("could not write periodic data [%s]", err.Error())
				if monitoringClient == c {
					monitoringClient = nil
				}
				return
			}
		case _ = <-c.close:
			c.log.Debugf("closing web-scokets")
			if monitoringClient == c {
				monitoringClient = nil
			}
			return
		}
	}
}

func serveWs(log slf.Logger) func(http.ResponseWriter, *http.Request) {
	return (func(w http.ResponseWriter, r *http.Request) {

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

		log.Debug("client: initdb")
		db, err := database.InitDB(globalOpt.DatabaseConfig)
		if err != nil {
			log.Fatalf("Could not init DB: [%s]", err.Error())
		}
		defer func() {
			log.Debug("Closing db")
			db.Close()
		}()
		log.Debug("client: initdb finished")

		c := &client{
			sendCh: make(chan []byte, maxMessageSize),
			ws:     ws,
			log:    slf.WithContext("client=" + ws.RemoteAddr().String()),
			db:     db,
		}

		go c.process()
		c.read()

		/*	if err := db.Close(); err != nil {
			log.Errorf("could not close db: %s", err.Error())
		}*/
	})
}
