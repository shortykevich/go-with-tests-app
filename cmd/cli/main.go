package main

import (
	"log"
	"net/http"

	"github.com/shortykevich/go-with-tests-app/db/inmem"
	"github.com/shortykevich/go-with-tests-app/webserver"
)

const port = ":5000"

func main() {
	handler := webserver.NewPlayersScoreServer(inmem.NewInMemoryStorage())

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
