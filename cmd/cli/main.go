package main

import (
	"fmt"
	"log"
	"os"

	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
	"github.com/shortykevich/go-with-tests-app/poker"
)

var dummyAlerter = &poker.SpyBlindAlerter{}

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker!")
	fmt.Println(`Type "{Name} wins" to record a win`)

	storage, close, err := fss.FileSystemStorageFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), storage)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
