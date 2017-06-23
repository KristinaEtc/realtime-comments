package main

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "test_comments"
)

func initDB() (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
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

func insertData(db *sql.DB) error {

	log.Debugf("# Inserting values")
	log.Debugf("")
	log.Debugf("Adding to db\n")

	sqlStatement := `  
INSERT INTO comment (user_name, comment, video_id, video_timestamp, calendar_timestamp)  
VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(sqlStatement, "vasia", "your test comment", 66, 4, time.Now())
	if err != nil {
		return fmt.Errorf("ading data to database: %s", err.Error())
	}
	return nil
}
