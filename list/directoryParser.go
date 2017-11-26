package list

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
)

const (
	watchDirFile = "watchDir"
	serveDirFile = "serveDir"
)

var songs []Song
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

		songs = append(songs, Song{m.Artist(), m.Title(), m.Album(), serveDir + path[len(watchDir):]})
	}
	return nil
}

//GetSongs returns Song objects based on what it finds in the watchedDir
func GetSongs() ([]Song, error) {
	songs = nil
	filepath.Walk(watchDir, visit)

	return songs, nil
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
