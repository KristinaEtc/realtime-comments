package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func initDB() (*sql.DB, error) {
	sqlOpenStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		globalOpt.DataBaseConfig.User, globalOpt.DataBaseConfig.Password, globalOpt.DataBaseConfig.NameDB, globalOpt.DataBaseConfig.Host)
	log.Info(sqlOpenStr)
	db, err := sql.Open("postgres", sqlOpenStr)
	if err != nil {
		log.Errorf("cannot opet database: %s", err.Error())
	}
	return db, err
}

//	checkErr(err)
//	defer db.Close()

func getData(db *sql.DB) error {
	rows, err := db.Query("SELECT * FROM comment")
	if err != nil {
		return fmt.Errorf("scan inserting data: %s", err.Error())
	}
	for rows.Next() {

		var (
			user_name          string
			comment            string
			video_id           int
			video_timestamp    sql.NullInt64
			calendar_timestamp time.Time
		)

		var i int
		err = rows.Scan(&user_name, &comment, &video_id, &video_timestamp, &calendar_timestamp)
		if err != nil {
			log.Errorf("scan inserting data: %s", err.Error())
			return fmt.Errorf("scan inserting data: %s", err.Error())
		}
		if video_timestamp.Valid {
			i = int(video_timestamp.Int64)
		} else {
			i = 0
		}
		log.Debugf("user_name|comment|video_id| video_timestamp| calendar_timestamp ")
		log.Debugf("%s\t | %s\t | %d  \t  | %d\t\t   | %v\n", user_name, comment, video_id, i, calendar_timestamp)
	}
	return nil
}

func insertData(db *sql.DB, data []byte, t time.Time) error {

	log.Debug("Adding to db")

	sqlStatement := `  
INSERT INTO comment (user_name, comment, video_id, video_timestamp, calendar_timestamp)  
VALUES ($1, $2, $3, $4, $5)`

	s := string(data)
	_, err := db.Query(sqlStatement, "vasia", s[:(len(s)-600)], 66, 4, t)
	if err != nil {
		log.Errorf("ading data to database: %s", err.Error())
		return fmt.Errorf("ading data to database: %s", err.Error())
	}
	return nil
}
