package main

import (
	"log"
	"net/http"

	"github.com/jackloughran/stream-yourself/art"
	"github.com/jackloughran/stream-yourself/list"
)

const (
	baseURL = "/api"
	listURL = baseURL + "/list"
	artURL  = baseURL + "/art"

	port = ":48001"
)

func main() {
	http.HandleFunc(listURL, list.Handler)
	http.HandleFunc(artURL, art.Handler)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
