package main

import _ "github.com/KristinaEtc/slflog"

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/KristinaEtc/config"
	"github.com/KristinaEtc/realtime-comments/database"
	"github.com/ventu-io/slf"
)

//DataBaseConf is a part of config with databse settings
type DataBaseConf struct {
	Type     string
	User     string
	Password string
	NameDB   string
	Host     string
}

//ServerConf is a part of config with server settings
type ServerConf struct {
	Address              string
	MessageSendPeriod    int
	FileWithTextData     string
	WriteTestDataTimeout int
	Broadcast            bool
	MonitoringMessage    string
}

// ConfFile is a file with all program options
type ConfFile struct {
	Name           string
	DataBaseConfig DataBaseConf
	ServerConfig   ServerConf
	/*WriteWait      time.Duration
	PongWait       time.Duration
	PingPeriod     time.Duration
	MaxMessageSize int
	*/
}

var globalOpt = ConfFile{
	Name: "WS-test",
	ServerConfig: ServerConf{
		Address:              "localhost:7777",
		MessageSendPeriod:    5,
		FileWithTextData:     "test-data2",
		WriteTestDataTimeout: 5,
		Broadcast:            false,
		MonitoringMessage:    "monitoring",
	},
	DataBaseConfig: DataBaseConf{
		Type:     "mock",
		User:     "guest",
		Password: "guest",
		NameDB:   "test",
		Host:     "localhost:5432",
	},
}

var data []byte
var db database.DataBase

func parseFileWithTextData() error {
	var err error
	data, err = ioutil.ReadFile(globalOpt.ServerConfig.FileWithTextData)
	if err != nil {
		return fmt.Errorf("No file with [%s] test data. Exiting", globalOpt.ServerConfig.FileWithTextData)
	}
	return nil
}

func main() {

	var log = slf.WithContext("realtime-comments")

	config.ReadGlobalConfig(&globalOpt, "WS-options")
	log.Errorf("-------------------------------------------")
	log.Infof("Running with next configuration: %+v", globalOpt)

	var err error
	db, err = database.InitDB(globalOpt.DataBaseConfig.Type)
	if err != nil {
		log.Panicf("Could not init DB: %s", err.Error())
	}
	defer db.Close()

	err = parseFileWithTextData()
	if err != nil {
		log.Errorf("No file with [%s] test data. Exiting", globalOpt.ServerConfig.FileWithTextData)
		os.Exit(1)
	}
	log.Debugf("Server will send data from a file [%s]", globalOpt.ServerConfig.FileWithTextData)

	// TODO: check if directory with html-stuff exists!
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", serveWs(log))
	log.Debug(globalOpt.ServerConfig.Address)
	err = http.ListenAndServe(globalOpt.ServerConfig.Address, nil)
	if err != nil {
		log.Errorf("Listen&Serve: [%s]", err.Error())
	}
}
