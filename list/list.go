package list

import (
	"encoding/json"
	"log"
	"net/http"
)

//Song represents basic information about songs found in the directory
type Song struct {
	ID       string `json:"id"`
	Artist   string `json:"artist"`
	Title    string `json:"title"`
	Album    string `json:"album"`
	CoverArt string `json:"coverArt"`
	Loc      string `json:"loc"`
}

//Handler handles requests to the /list endpoint which lists songs
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request to /list. [uri: %s]", r.RequestURI)

	err := Parse()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(Songs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
