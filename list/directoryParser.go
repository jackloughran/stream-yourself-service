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
	configFile = "config"
)

var songs []Song
var watchedDir = getWatchDir()

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

		songs = append(songs, Song{m.Artist(), m.Title(), m.Album(), path})
	}
	return nil
}

//GetSongs returns Song objects based on what it finds in the watchedDir
func GetSongs() ([]Song, error) {
	songs = nil
	filepath.Walk(watchedDir, visit)

	return songs, nil
}

func getWatchDir() string {
	dat, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error reading config file: " + err.Error())
	}

	return strings.Trim(string(dat), "\n")
}
