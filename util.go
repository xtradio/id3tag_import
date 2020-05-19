package main

import (
	"database/sql"
	"fmt"
	"os"
)

func getEnv(envKey string) (envValue string, err error) {

	envValue, ok := os.LookupEnv(envKey)
	if ok != true {
		err = fmt.Errorf("please set %s environment variable", envKey)
		return
	}

	return envValue, nil
}

func dbConnection() (*sql.DB, error) {
	username, err := getEnv("MYSQL_USERNAME")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("MYSQL_PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := getEnv("MYSQL_HOST")
	if err != nil {
		return nil, err
	}

	database, err := getEnv("MYSQL_DATABASE")
	if err != nil {
		return nil, err
	}

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", username, password, host, database)

	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getSongsFromDB(db *sql.DB) ([]SongDetails, error) {
	var l []SongDetails

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	rows, err := db.Query("SELECT id, filename, artist, title, album, lenght, share, url, image, playlist FROM details ORDER BY id DESC")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s SongDetails

		rows.Scan(&s.ID, &s.Filename, &s.Artist, &s.Title, &s.Album, &s.Length, &s.Share, &s.URL, &s.Image, &s.Playlist)

		l = append(l, s)

	}

	defer rows.Close()

	return l, nil
}

func doesExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
