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
	db, closeDb := fss.InitDB(dbFileName)
	defer closeDb()

	handler := webserver.NewPlayersScoreServer(fss.NewFSPlayerStorage(db))

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
