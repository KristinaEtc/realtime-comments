package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // for postgres
)

// PostgresDB is a struct which wrapps sql.DB entity.
// It implements DataBase interface.
type PostgresDB struct {
	db   *sql.DB
	conf Conf
}

func initPostgresDB(conf Conf) (Database, error) {

	fmt.Print("initPostgresDB")

	sqlOpenStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s  sslmode=disable",
		conf.User, conf.Password, conf.NameDB, conf.Host)

	db, err := sql.Open("postgres", sqlOpenStr)
	if err != nil {
		return nil, fmt.Errorf("cannot opet database: %s", err.Error())
	}

	//db.SetConnMaxLifetime(time.Second * 10)
	//db.SetMaxIdleConns(500)
	db.SetMaxOpenConns(2000)

	postgresDB := &PostgresDB{
		db:   db,
		conf: conf,
	}

	fmt.Printf("%+v\n", db)
	//os.Exit(1)

	var tmp int
	err = db.QueryRow("SELECT 1").Scan(&tmp)
	if err != nil {
		log2.Errorf("db.QueryRow %s", err.Error())
		return nil, err
	}
	return postgresDB, nil
}

// GetData returns data from special table.
func (p *PostgresDB) GetData() ([]byte, error) {
	/*rows, err := p.db.Query("SELECT * FROM ? ORDER BY id DESC LIMIT 10;", p.conf.Table)
	if err != nil {
		return nil, fmt.Errorf("scan inserting data: %s", err.Error())
	}

	var (
		user_name          string
		comment            string
		video_id           int
		video_timestamp    sql.NullInt64
		calendar_timestamp time.Time
	)
	var i int

	for rows.Next() {

		err = rows.Scan(&user_name, &comment, &video_id, &video_timestamp, &calendar_timestamp)
		if err != nil {
			log.Errorf("scan inserting data: %s", err.Error())
			return nil, fmt.Errorf("scan inserting data: %s", err.Error())
		}
		if video_timestamp.Valid {
			i = int(video_timestamp.Int64)
		} else {
			i = 0
		}
		//TODO: return json
		log.Debugf("user_name|comment|video_id| video_timestamp| calendar_timestamp ")
		log.Debugf("%s\t | %s\t | %d  \t  | %d\t\t   | %v\n", user_name, comment, video_id, i, calendar_timestamp)

	}*/

	const (
		user_name       = "vasia"
		comment         = "my comment"
		video_id        = 14
		video_timestamp = 15
	)
	calendar_timestamp := time.Now()
	var i int
	return []byte(fmt.Sprintf("%s\t | %s\t | %d  \t  | %d\t\t   | %v\n", user_name, comment, video_id, i, calendar_timestamp)), nil
}

// InsertData inserts data to special table.
//func (p *PostgresDB) InsertData(db *sql.DB, data []byte, t time.Time, log slf.Logger) error {
func (p *PostgresDB) InsertData(data []byte, currTime time.Time) error {

	fmt.Println("Adding to db")

	sqlStatement := `  
INSERT INTO comment (user_name, comment, video_id, video_timestamp, calendar_timestamp)  
VALUES ($1, $2, $3, $4, $5)`

	s := string(data)
	_, err := p.db.Exec(sqlStatement, "vasia", s[:(len(s)-600)], 66, 4, currTime)
	if err != nil {
		fmt.Println("ading data to database: ", err)
		return fmt.Errorf("ading data to database: %s", err.Error())
	}
	return nil
}

// Close is a method which closes postgres DB.
func (p *PostgresDB) Close() error {

	fmt.Println("Closing postgres DB")
	return p.Close()
}
