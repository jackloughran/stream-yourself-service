package art

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dhowden/tag"
	"github.com/jackloughran/stream-yourself/list"
)

var fileName string
var songDir string
var picture tag.Picture
var storedPic storedPicture

type storedPicture struct {
	Type string
	Data []byte
}

// Handler handles requests to get album art
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request to /art. [uri: %s]", r.RequestURI)

	fileNameSlice := r.URL.Query()["fileName"]
	if fileNameSlice == nil {
		http.Error(w, "expecting fileName query parameter", http.StatusInternalServerError)
		log.Printf("expecting fileName query parameter")
	}

	fileName = fileNameSlice[0]

	picture = tag.Picture{}
	filepath.Walk(list.WatchDir, visit)

	if picture.Data != nil {
		w.Header().Set("Content-Type", picture.MIMEType)
		w.Header().Set("Content-Length", strconv.Itoa(len(picture.Data)))
		w.Write(picture.Data)
	} else {
		log.Printf("Checking %s for a .jpg or a .png file", songDir)
		storedPic = storedPicture{}
		filepath.Walk(songDir, visit2)
		if storedPic.Data != nil {
			w.Header().Set("Content-Type", storedPic.Type)
			w.Header().Set("Content-Length", strconv.Itoa(len(storedPic.Data)))
			w.Write(storedPic.Data)
		} else {
			http.Error(w, "Unable to find picture.", http.StatusNotFound)
		}
	}

}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && f.Name() == fileName {
		// first store the directory in case we cant find the picture on the flac file
		songDir = path[:len(path)-len(f.Name())]

		file, err := os.Open(path)
		if err != nil {
			log.Printf("Error trying to read file for art: %v", err)
			return err
		}

		m, err := tag.ReadFrom(file)
		if err != nil {
			log.Printf("Error during tag in art: %v", err)
			return err
		}

		if m != nil && m.Picture() != nil {
			picture = *m.Picture()
		} else {
			log.Printf("Unable to find picture for %s", m.Title())
		}
	}

	return nil
}

func visit2(path string, f os.FileInfo, err error) error {
	if !f.IsDir() {
		fName := f.Name()
		if fName[len(fName)-4:] == ".jpg" || fName[len(fName)-4:] == ".png" {
			log.Printf("Found stored cover file: %s", path)
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				log.Print("Error opening stored art file")
				return err
			}

			if fName[len(fName)-4:] == ".jpg" {
				storedPic.Type = "image/jpeg"
			} else {
				storedPic.Type = "image/png"
			}
			storedPic.Data = dat
		}
	}

	return nil
}
