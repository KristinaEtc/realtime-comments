package main

import (
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/KristinaEtc/config"
	"github.com/KristinaEtc/realtime-comments/database"
	_ "github.com/KristinaEtc/slflog"
	"github.com/gorilla/websocket"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("ws-client")

// ConfFile is a file with all program options
type ConfFile struct {
	Name            string
	Address         string
	ConnectionCount int
	MessageCount    int
	PeriodToSend    int
}

var globalOpt = ConfFile{
	Name:            "WS-client",
	Address:         "ws://localhost:7778/ws",
	MessageCount:    10,
	ConnectionCount: 100,
	PeriodToSend:    5,
}

var connectedCount int

func runTest(wg *sync.WaitGroup, id int) {

	defer wg.Done()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(globalOpt.Address, nil)
	if err != nil {
		log.Panicf("dial: %s", err.Error())
	}
	connectedCount++
	log.Infof("connected %d", connectedCount)
	defer c.Close()

	/*u, err := uuid.NewV4()
	if err != nil {
		log.Panicf("Could not generate uuid: %s", err.Error())
	}*/

	var commentData = database.Comment{
		CommentBody:       strconv.Itoa(id) + " HelloWorld",
		Username:          strconv.Itoa(id),
		VideoID:           6,
		VideoTimestamp:    6,
		CalendarTimestamp: time.Now(),
	}

	var commentJSON []byte
	if commentJSON, err = json.Marshal(&commentData); err != nil {
		log.Errorf("could serialize json [%+v]: %s", commentJSON, err.Error())
		return
	}

	log.Debugf("sending %s", string(commentJSON))

	err = c.WriteMessage(websocket.TextMessage, commentJSON)
	if err != nil {
		log.WithField("id", id).Errorf("write: %s", err.Error())
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(globalOpt.PeriodToSend))
	defer ticker.Stop()
	//done := make(chan struct{})
	done := make(chan bool)

	go func() {

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.WithField("id", id).Errorf("read: %s", err.Error())
				done <- true
				return
			}
			log.Debugf("Got %s", string(message))
			log.WithField("id", id).Debugf("recv: [%s]", message)
		}
	}()

	for {
		select {
		case _ = <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, commentJSON)
			if err != nil {
				log.Errorf("write: %s", err.Error())
				return
			}
		case _ = <-done:
			{
				log.WithField("id", id).Debug("closing")
				return
			}
		case <-interrupt:
			log.Warn("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.WithField("id", id).Errorf("write close: %s", err.Error())
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}

func main() {

	log.Error("------------------------------------------------")
	config.ReadGlobalConfig(&globalOpt, "WS-options")
	log.Infof("%v", globalOpt)
	time.Sleep(time.Second * 2)

	var wg sync.WaitGroup
	wg.Add(globalOpt.ConnectionCount)
	for i := 0; i < globalOpt.ConnectionCount; i++ {
		go runTest(&wg, i)
	}

	wg.Wait()
	log.Infof("Connections=[%d]", connectedCount)
}
