package main

import _ "github.com/KristinaEtc/slflog"

import (
	"io/ioutil"
	"net/http"
	"os"

	"database/sql"

	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("realtime-comments")

//DataBaseConf is a part of config with databse settings
type DataBaseConf struct {
	User     string
	Password string
	NameDB   string
	Host     string
}

// ConfFile is a file with all program options
type ConfFile struct {
	Name                 string
	Address              string
	MessageSendPeriod    int
	FileWithTextData     string
	WriteTestDataTimeout int
	Broadcast            bool
	MonitoringMessage    string
	DataBaseConfig       DataBaseConf
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
	Broadcast:            false,
	MonitoringMessage:    "monitoring",
	DataBaseConfig: DataBaseConf{
		User:     "guest",
		Password: "guest",
		NameDB:   "test",
		Host:     "localhost:5432",
	},
	/*	WriteWait:      10 * time.Second,
		PongWait:       10 * time.Second,
		PingPeriod:     5 * time.Second,
		MaxMessageSize: 1024 * 1024,
	*/
}

var data []byte
var db *sql.DB

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
	log.Errorf("-------------------------------------------")
	log.Infof("Running with next configuration: %+v", globalOpt)

	var err error
	db, err = initDB()
	if err != nil {
		log.Panicf("Could not init DB: %s", err.Error())
	}
	defer db.Close()

	parseFileWithTextData()

	// TODO: check if directory with html-stuff exists!
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs)
	log.Debug(globalOpt.Address)
	err = http.ListenAndServe(globalOpt.Address, nil)
	if err != nil {
		log.Errorf("Listen&Serve: [%s]", err.Error())
	}
}
