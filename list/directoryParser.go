package list

import (
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dhowden/tag"
)

const (
	watchDirFile = "watchDir"
	serveDirFile = "serveDir"
)

//Song represents basic information about songs found in the directory
type Song struct {
	ID     string `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
	Album  string `json:"album"`
	ArtLoc string `json:"artLoc"`
	Loc    string `json:"loc"`
}

// Songs lists all the songs, artists, and albums
var Songs []Song

// WatchDir holds the directory where the music is stored
var WatchDir = getWatchDir()
var serveDir = getServeDir()

func visit(path string, f os.FileInfo, err error) error {
	fileName := f.Name()
	if !f.IsDir() && fileName[len(fileName)-5:] == ".flac" {
		file, err := os.Open(path)
		if err != nil {
			log.Printf("Error during visit open: %v", err)
			return err
		}

		m, err := tag.ReadFrom(file)

		if err != nil {
			log.Printf("Error during visit tag: %v", err)
			return err
		}

		Songs = append(Songs, Song{strconv.Itoa(len(Songs)), m.Artist(), m.Title(), m.Album(), "http://138.197.172.114:48001/api/art?fileName=" + url.QueryEscape(fileName), serveDir + path[39:]})
	}

	return nil
}

// Parse initiates a walkthrough of the directory to find songs and pictures
func Parse() error {
	Songs = nil
	filepath.Walk(WatchDir, visit)

	return nil
}

func getWatchDir() string {
	dat, err := ioutil.ReadFile(watchDirFile)
	if err != nil {
		log.Fatal("Error reading config file: " + err.Error())
	}

	return strings.Trim(string(dat), "\n")
}

func getServeDir() string {
	dat, err := ioutil.ReadFile(serveDirFile)
	if err != nil {
		log.Fatal("Error reading config file: " + err.Error())
	}

	return strings.Trim(string(dat), "\n")
}
