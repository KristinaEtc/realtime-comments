package main

import _ "github.com/KristinaEtc/slflog"

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("Web-socket-test")

// ConfFile is a file with all program options
type ConfFile struct {
	Name                 string
	Address              string
	MessageSendPeriod    int
	FileWithTextData     string
	WriteTestDataTimeout int
	/*WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int
	*/
}

var globalOpt = ConfFile{
	Name:                 "WS-test",
	Address:              "localhost:7777",
	MessageSendPeriod:    5,
	FileWithTextData:     "test-data2",
	WriteTestDataTimeout: 5,
	/*	WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     5 * time.Second,
		MaxMessageSize: 1024 * 1024,
	*/
}

var data []byte

func parseFileWithTextData() {
	var err error
	data, err = ioutil.ReadFile(globalOpt.FileWithTextData)
	if err != nil {
		log.Errorf("No file with [%s] test data. Exiting", globalOpt.FileWithTextData)
		os.Exit(1)
	}
	log.Debugf("Server will send data from a file [%s]", globalOpt.FileWithTextData)
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "WS-options")
	log.Infof("Running with next configuration: %+v", globalOpt)

	parseFileWithTextData()

	go h.run()

	// TODO: check if directory with html-stuff exists!
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	log.Debug(globalOpt.Address)
	err := http.ListenAndServe(globalOpt.Address, nil)
	if err != nil {
		log.Errorf("Listen&Serve: [%s]", err.Error())
	}
}
