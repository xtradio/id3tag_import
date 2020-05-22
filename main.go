package main

import (
	"fmt"
	"log"
	"os"

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

	args := os.Args[1:]

	if len(args) > 1 {
		panic("Please only specify one argument.")
	}

	if args[0] == "synctags" {
		syncTags()
		// fmt.Println("Run synctags()")
	} else if args[0] == "fixduration" {
		log.Println("Running duration.")
		fixDuration()
	} else {
		fmt.Println("Please use synctags or fixduration as command line arguments.")
	}

}
