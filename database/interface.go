package database

import (
	"strings"
	"time"

	"github.com/ventu-io/slf"
)

var pwdCurr = "database"
var log = slf.WithContext(pwdCurr)

// DataBase is an interface for different databases. Moreover,
// it's useful, when no DB was set (in this case, MockDB would be used).
type DataBase interface {
	GetData() error
	InsertData([]byte, time.Time) error
	Close() error
}

// InitDB inits database according database parameter; if necessary BD wasn't found, function inits mock
func InitDB(nameDB string) (DataBase, error) {

	switch strings.ToLower(nameDB) {
	case "postgres":
		return initPostgresDB()
	default:
		return initMockDB()
	}
}
