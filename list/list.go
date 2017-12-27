package list

import (
	"encoding/json"
	"log"
	"net/http"
)

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
