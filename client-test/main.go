package main

import _ "github.com/KristinaEtc/slflog"
import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/KristinaEtc/config"

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

func runTest(wg *sync.WaitGroup) {

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

	//id := strconv.Itoa(connectedCount) + "" + u.String()
	id := strconv.Itoa(connectedCount) + " HelloWorld "

	//data := append([]byte(time.Now().Format("2006-01-02/15:04:05 ")), []byte(u.String())...)
	data := append([]byte(time.Now().Format("2006-01-02/15:04:05 ")), []byte(strconv.Itoa(connectedCount))...)
	err = c.WriteMessage(websocket.TextMessage, data)
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
			log.WithField("id", id).Debugf("recv: [%s]", message)
		}
	}()

	for {
		select {
		case _ = <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(id))
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
		go runTest(&wg)
	}

	wg.Wait()
	log.Infof("Connections=[%d]", connectedCount)
}
