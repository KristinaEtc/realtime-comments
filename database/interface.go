package database

import (
	"strings"
)

//Conf is a part of config with databse settings
type Conf struct {
	Type                  string
	User                  string
	Password              string
	NameDB                string
	Host                  string
	Table                 string
	NumOfSelectedComments int
}

// Database is an interface for different databases. Moreover,
// it's useful, when no DB was set (in this case, MockDB would be used).
type Database interface {
	GetLastComments(int64, int) ([]Comment, error)
	InsertData(Comment) error
	Close() error
}

// InitDB inits database according database parameter; if necessary BD wasn't found, function inits mock
func InitDB(config Conf) (Database, error) {
	log.Debugf("%v", config)

	switch strings.ToLower(config.Type) {
	case "postgres":
		log.Debug("InitDB: postgres")
		return initPostgresDB(config)
	default:
		return initMockDB()
	}
}
