package main

import (
	"net/http"
	"github.com/jackloughran/stream-yourself/list"
	"log"
)

const (
	baseURL = "/api"
	listURL = baseURL + "/list"

	port = ":48001"
)

func main() {
	http.HandleFunc(listURL, list.Handler)

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
