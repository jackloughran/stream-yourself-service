package list

import (
	"io/ioutil"
	"log"
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

// Songs lists all the songs, artists, and albums
var Songs []Song
var watchDir = getWatchDir()
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

		var coverArt string

		folder := serveDir + path[len(watchDir):len(path)-len(f.Name())]
		if art, err := filepath.Glob(folder + "*.jpg"); art != nil {
			coverArt = art[0]
		} else if err != nil {
			log.Printf("error searching for cover art: %v", err)
		}

		if art, err := filepath.Glob(folder + "*.png"); art != nil {
			coverArt = art[0]
		} else if err != nil {
			log.Printf("error searching for cover art: %v", err)
		}

		Songs = append(Songs, Song{strconv.Itoa(len(Songs)), m.Artist(), m.Title(), m.Album(), coverArt, serveDir + path[39:]})
	}

	return nil
}

// Parse initiates a walkthrough of the directory to find songs and pictures
func Parse() error {
	Songs = nil
	filepath.Walk(watchDir, visit)

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
