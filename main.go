package main

import (
	"log"
	"net/http"
)

const port = ":5000"

func main() {
	handler := &PlayersScoreServer{Storage: NewInMemoryStorage()}

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
