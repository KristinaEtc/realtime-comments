package main

import _ "github.com/KristinaEtc/slflog"

import (
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

	c, _, err := websocket.DefaultDialer.Dial(globalOpt.Address, nil)
	if err != nil {
		log.Panicf("dial: %s", err.Error())
	}
	connectedCount++
	log.Infof("connected %d", connectedCount)
	defer c.Close()

	done := make(chan struct{})
	/*u, err := uuid.NewV4()
	if err != nil {
		log.Panicf("Could not generate uuid: %s", err.Error())
	}*/

	id := strconv.Itoa(connectedCount) //+ "" + u.String()

	//data := append([]byte(time.Now().Format("2006-01-02/15:04:05 ")), []byte(u.String())...)
	data := append([]byte(time.Now().Format("2006-01-02/15:04:05 ")), []byte(strconv.Itoa(connectedCount))...)
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.WithField("id", id).Errorf("write: %s", err.Error())
		return
	}

	for {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.WithField("id", id).Errorf("read: %s", err.Error())
				return
			}
			log.WithField("id", id).Debugf("recv: [%s]", message)
		}
	}
}

func main() {

	log.Error("------------------------------------------------")
	config.ReadGlobalConfig(&globalOpt, "WS-options")
	log.Infof("%v", globalOpt)
	time.Sleep(time.Second * 5)

	var wg sync.WaitGroup
	wg.Add(globalOpt.ConnectionCount)
	for i := 0; i < globalOpt.ConnectionCount; i++ {
		go runTest(&wg)
	}

	wg.Wait()
	log.Infof("Connections=[%d]", connectedCount)
}
