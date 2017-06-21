package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

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

func runTest(wg *sync.WaitGroup) {

	defer wg.Done()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(globalOpt.Address, nil)
	if err != nil {
		log.Panicf("dial: %s", err.Error())
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Errorf("read: %s", err.Error())
				return
			}
			log.Debugf("recv: [%s]", message)
		}
	}()

	/*
		ticker := time.NewTicker(time.Second * time.Duration(globalOpt.PeriodToSend))
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(t.Format("2006-01-02/15:04:05 ")))
				if err != nil {
					log.Errorf("write: %s", err.Error())
					return
				}
			case <-interrupt:
				log.Warn("interrupt")
				// To cleanly close a connection, a client should send a close
				// frame and wait for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Errorf("write close: %s", err.Error())
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
	*/

	for i := 0; i < globalOpt.MessageCount; i++ {
		err := c.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02/15:04:05 ")))
		if err != nil {
			log.Errorf("write: %s", err.Error())
			return
		}
	}
	return
}

func main() {

	var wg sync.WaitGroup
	wg.Add(globalOpt.ConnectionCount)
	for i := 0; i < globalOpt.ConnectionCount; i++ {
		go runTest(&wg)
	}

	wg.Wait()
}
