package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/shortykevich/go-with-tests-app/cli"
	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
)

var dummyAlerter = &poker.SpyBlindAlerter{}

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker!")
	fmt.Println(`Type "{Name} wins" to record a win`)

	store, close, err := fss.FileSystemStorageFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewCLI(store, os.Stdin, os.Stdout, dummyAlerter)
	game.PlayPoker()
}
