package main

import (
	"fmt"
	"log"

	"github.com/bogem/id3v2"
	_ "github.com/go-sql-driver/mysql"
)

// SongDetails to output details of the songs to json
type SongDetails struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Show     string `json:"show"`
	Image    string `json:"image"`
	Filename string `json:"filename"`
	Album    string `json:"album"`
	Length   string `json:"lenght"`
	Share    string `json:"share"`
	URL      string `json:"url"`
	Playlist string `json:"playlist"`
}

func main() {

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
		if !doesExist(path) {
			// fmt.Println(path)
			continue
		}
		fmt.Println(path)

		err := syncTags(song, path)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func syncTags(song SongDetails, path string) error {
	tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
	if err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
		return err
	}
	defer tag.Close()

	tagArtist, tagTitle := tag.Artist(), tag.Title()

	if tagArtist != song.Artist {
		tag.SetArtist(song.Artist)
	}

	if tagTitle != song.Title {
		tag.SetTitle(song.Title)
	}

	link := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "link",
		Text:        song.Share,
	}
	tag.AddCommentFrame(link)

	imageLink := id3v2.CommentFrame{
		Encoding:    id3v2.EncodingUTF8,
		Language:    "eng",
		Description: "imageLink",
		Text:        song.Image,
	}
	tag.AddCommentFrame(imageLink)

	// Write tag to file.mp3.
	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
		return err
	}

	return nil
}
