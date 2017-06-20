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
	Name string
	/*WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int*/
}

var globalOpt = ConfFile{
	Name: "WS-test",
	/*	WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     5 * time.Second,
		MaxMessageSize: 1024 * 1024,*/
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "WS-options")
	log.Infof("%+v", globalOpt)

	go h.run()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("Listen&Serve: %s", err.Error())
	}
}
