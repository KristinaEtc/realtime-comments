package main

import _ "github.com/KristinaEtc/slflog"

import (
	"net/http"

	"github.com/KristinaEtc/config"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("WS-test")

// ConfFile is a file with all program options
type ConfFile struct {
	Name              string
	Address           string
	MessageSendPeriod int
	/*WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int*/
}

var globalOpt = ConfFile{
	Name:              "WS-test",
	Address:           "localhost:7777",
	MessageSendPeriod: 5,
	/*	WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     5 * time.Second,
		MaxMessageSize: 1024 * 1024,*/
}

var data []byte

func main() {

	config.ReadGlobalConfig(&globalOpt, "WS-options")
	log.Infof("%+v", globalOpt)

	var err error
	//	data, err = ioutil.ReadFile("test-data2")
	//	if err != nil {
	//		log.Errorf("No file with test data. Exiting")
	//		os.Exit(1)

	//	} else {
	//		log.Debug("all ok")
	//	}

	go h.run()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	log.Debug(globalOpt.Address)
	err = http.ListenAndServe(globalOpt.Address, nil)
	if err != nil {
		log.Errorf("Listen&Serve: %s", err.Error())
	}
}
