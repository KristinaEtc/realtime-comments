package database

import (
	"strings"
	"time"

	"github.com/ventu-io/slf"
)

var pwdCurr = "database"
var log2 = slf.WithContext(pwdCurr)

//Conf is a part of config with databse settings
type Conf struct {
	Type     string
	User     string
	Password string
	NameDB   string
	Host     string
	Table    string
}

// Database is an interface for different databases. Moreover,
// it's useful, when no DB was set (in this case, MockDB would be used).
type Database interface {
	GetData() ([]byte, error)
	InsertData([]byte, time.Time) error
	Close() error
}

// InitDB inits database according database parameter; if necessary BD wasn't found, function inits mock
func InitDB(config Conf) (Database, error) {
	log2.Debugf("%v", config)

	switch strings.ToLower(config.Type) {
	case "postgres":
		log2.Debug("InitDB: postgres")
		return initPostgresDB(config)
	default:
		return initMockDB()
	}
}
