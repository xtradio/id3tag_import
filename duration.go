package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tcolgate/mp3"
)

func fixDuration() {
	var v struct {
		Data []SongDetails `json:"data"`
	}

	db, err := dbConnection()
	if err != nil {
		log.Println("Error opening connection to DB: ", err)
		return
	}
	defer db.Close()

	v.Data, err = getSongsFromDB(db)
	if err != nil {
		log.Println("Error getting items from database: ", err)
		return
	}

	localPath, err := getEnv("MUSIC_LOCAL_PATH")
	if err != nil {
		return
	}

	for _, song := range v.Data {
		path := fmt.Sprintf("%s%s", localPath, song.Filename)
		id := song.ID
		if !doesExist(path) {
			log.Println("File does not exist: ", path)
			break
		}
		log.Println(path)

		songDuration, err := duration(path)
		if err != nil {
			log.Println("Can't get duration of song: ", path, err)
			return
		}

		err = updateSong(db, id, songDuration)
		if err != nil {
			log.Println("Failed to update song in DB: ", err)
			return
		}

	}
}

func duration(file string) (float64, error) {
	t := 0.0

	r, err := os.Open(file)
	if err != nil {
		return t, err
	}
	defer r.Close()

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return t, err
		}

		t = t + f.Duration().Seconds()
	}

	return t, nil

}

func updateSong(db *sql.DB, id int64, duration float64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("UPDATE details SET lenght = ? WHERE id = ?", duration, id); err != nil {
		return err
	}

	return nil
}
