package database

import (
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

// Comment is a struct for comment data
type Comment struct {
	ID                int
	Username          string
	CommentBody       string
	VideoID           int64
	VideoTimestamp    int64
	CalendarTimestamp time.Time
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

// GetLastComments returns last 10 comments.
func (p *PostgresDB) GetLastComments() ([]Comment, error) {

	log.Debug("GetLastComments")

	// TODO: !!! VERY IMPORTANT: ADD TABLE AS ARGUMENT
	//rows, err := p.db.Query(`SELECT * FROM $1 ORDER BY id DESC LIMIT 10;`, p.conf.Table)
	rows, err := p.db.Query(`SELECT * FROM comment ORDER BY id DESC LIMIT 10;`)
	if err != nil {
		return nil, fmt.Errorf("scan inserting data: %s", err.Error())
	}

	var comments = make([]Comment, 10)
	var videoTimestamp sql.NullInt64

	for rows.Next() {

		var comment = Comment{}
		err = rows.Scan(&comment.Username, &comment.CommentBody, &comment.VideoID, &videoTimestamp, &comment.CalendarTimestamp, &comment.ID)
		if err != nil {
			log.Errorf("scan inserting data: %s", err.Error())
			return nil, fmt.Errorf("scan inserting data: %s", err.Error())
		}
		if videoTimestamp.Valid {
			comment.VideoTimestamp = videoTimestamp.Int64
		} else {
			log.Errorf("Parse [%v] video_timestamp: not valid. Set a zero.", videoTimestamp)
			comment.VideoTimestamp = 0
		}
		//log.Debugf("user_name|comment|video_id| video_timestamp| calendar_timestamp ")
		//log.Debugf("%s\t | %s\t | %d  \t  | %d\t\t   | %v\n", user_name, comment, video_id, i, calendar_timestamp)

		comments = append(comments, comment)
	}
	return comments, nil
}

// InsertData inserts data to special table.
//func (p *PostgresDB) InsertData(db *sql.DB, data []byte, t time.Time, log slf.Logger) error {
func (p *PostgresDB) InsertData(commentData Comment) error {

	log.Debug("Adding to db")

	sqlStatement := `  
INSERT INTO comment (user_name, comment, video_id, video_timestamp, calendar_timestamp)  
VALUES ($1, $2, $3, $4, $5)`

	_, err := p.db.Exec(sqlStatement, commentData.Username, commentData.CommentBody, commentData.VideoID,
		commentData.VideoTimestamp, commentData.CalendarTimestamp)
	if err != nil {
		log.Errorf("ading data to database: %s ", err.Error())
		return fmt.Errorf("ading data to database: %s", err.Error())
	}
	return nil
}

// Close is a method which closes postgres DB.
func (p *PostgresDB) Close() error {

	fmt.Println("Closing postgres DB")
	return p.db.Close()
}
