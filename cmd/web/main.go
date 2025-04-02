package main

import (
	"log"
	"net/http"

	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
	"github.com/shortykevich/go-with-tests-app/webserver"
)

const (
	port       = ":5000"
	dbFileName = "game.db.json"
)

func main() {
	store, close, err := fss.FileSystemStorageFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	handler := webserver.NewPlayersScoreServer(store)

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
