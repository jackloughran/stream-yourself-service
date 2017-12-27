package art

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dhowden/tag"
	"github.com/jackloughran/stream-yourself/list"
)

var fileName string
var picture tag.Picture

// Handler handles requests to get album art
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request to /art. [uri: %s]", r.RequestURI)

	fileNameSlice := r.URL.Query()["fileName"]
	if fileNameSlice == nil {
		http.Error(w, "expecting fileName query parameter", http.StatusInternalServerError)
		log.Printf("expecting fileName query parameter")
	}

	fileName = fileNameSlice[0]

	filepath.Walk(list.WatchDir, visit)

	w.Header().Set("Content-Type", picture.MIMEType)
	w.Header().Set("Content-Length", strconv.Itoa(len(picture.Data)))
	w.Write(picture.Data)

}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && f.Name() == fileName {
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

		picture = *m.Picture()
	}

	return nil
}
