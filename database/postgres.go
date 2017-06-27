package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // for postgres
	"github.com/ventu-io/slf"
)

var pwdCurr = "database"
var log = slf.WithContext(pwdCurr)

// PostgresDB is a struct which wrapps sql.DB entity.
// It implements DataBase interface.
type PostgresDB struct {
	db *sql.DB
	//sync.Mutex
	conf Conf
}

func initPostgresDB(conf Conf) (Database, error) {
	fmt.Print("initPostgresDB")

	sqlOpenStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		conf.User, conf.Password, conf.NameDB, conf.Host)

	db, err := sql.Open("postgres", sqlOpenStr)
	if err != nil {
		return nil, fmt.Errorf("cannot opet database: %s", err.Error())
	}

	db.SetConnMaxLifetime(time.Second * 3)
	// db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(2000)

	postgresDB := &PostgresDB{
		db:   db,
		conf: conf,
	}

	fmt.Printf("%+v\n", db)
	//os.Exit(1)

	if err = db.Ping(); err != nil {
		log.Errorf("db.QueryRow %s", err.Error())
		return nil, err
	}
	return postgresDB, nil
}

// GetData returns data from special table.
func (p *PostgresDB) GetData() ([]byte, error) {

	log.Debug("Get data")
	//rows, err := p.db.Query(`SELECT * FROM $1 ORDER BY id DESC LIMIT 10;`, p.conf.Table)
	rows, err := p.db.Query(`SELECT * FROM comment ORDER BY id DESC LIMIT 10;`)
	if err != nil {
		return nil, fmt.Errorf("scan inserting data: %s", err.Error())
	}

	var buffer bytes.Buffer

	for rows.Next() {

		var (
			id                 int64
			user_name          string
			comment            string
			video_id           int
			video_timestamp    sql.NullInt64
			calendar_timestamp time.Time
		)
		var i int
		err = rows.Scan(&user_name, &comment, &video_id, &video_timestamp, &calendar_timestamp, &id)
		if err != nil {
			log.Errorf("scan inserting data: %s", err.Error())
			return nil, fmt.Errorf("scan inserting data: %s", err.Error())
		}
		if video_timestamp.Valid {
			i = int(video_timestamp.Int64)
		} else {
			i = 0
		}
		//log.Debugf("user_name|comment|video_id| video_timestamp| calendar_timestamp ")
		//log.Debugf("%s\t | %s\t | %d  \t  | %d\t\t   | %v\n", user_name, comment, video_id, i, calendar_timestamp)

		//TODO: return json
		buffer.WriteString(fmt.Sprintf("%s\t | %s\t | %d  \t  | %d\t\t   | %v\n", user_name, comment, video_id, i, calendar_timestamp))
	}
	return buffer.Bytes(), nil
}

// InsertData inserts data to special table.
//func (p *PostgresDB) InsertData(db *sql.DB, data []byte, t time.Time, log slf.Logger) error {
func (p *PostgresDB) InsertData(data []byte, currTime time.Time) error {

	log.Debug("Adding to db")

	sqlStatement := `  
INSERT INTO comment (user_name, comment, video_id, video_timestamp, calendar_timestamp)  
VALUES ($1, $2, $3, $4, $5)`

	// s := string(data)
	_, err := p.db.Exec(sqlStatement, "vasia", "comment", 66, 4, currTime)
	if err != nil {
		fmt.Println("ading data to database: ", err)
		return fmt.Errorf("ading data to database: %s", err.Error())
	}
	return nil
}

// Close is a method which closes postgres DB.
func (p *PostgresDB) Close() error {

	fmt.Println("Closing postgres DB")
	return p.db.Close()
}
