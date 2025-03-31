package main

import (
	"log"
	"net/http"
	"os"

	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
	"github.com/shortykevich/go-with-tests-app/webserver"
)

const (
	port       = ":5000"
	dbFileName = "game.db.json"
)

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Problem opening %s %v", dbFileName, err)
	}

	handler := webserver.NewPlayersScoreServer(fss.NewFSPlayerStorage(db))

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
