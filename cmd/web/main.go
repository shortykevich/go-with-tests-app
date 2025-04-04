package main

import (
	"log"
	"net/http"

	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
	"github.com/shortykevich/go-with-tests-app/poker"
	"github.com/shortykevich/go-with-tests-app/webserver"
)

const (
	port       = ":5000"
	dbFileName = "game.db.json"
)

func main() {
	storage, close, err := fss.FileSystemStorageFromFile(dbFileName)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}
	defer close()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), storage)

	handler, err := webserver.NewPlayersScoreServer(storage, game)
	if err != nil {
		log.Fatalf("problem creating player server %v", err)
	}

	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
