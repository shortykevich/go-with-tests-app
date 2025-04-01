package fss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/league"
)

type FileSystemPlayerStorage struct {
	mu     sync.Mutex
	Db     io.ReadWriteSeeker
	League league.League
}

// Function to initialize db (json) file
// Return *os.File and function to close it
func InitDB(path string) (*os.File, func()) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Problem opening %s %v", path, err)
	}
	// Just for the sake of my sanity. Please one less "if err != nil"
	stat, _ := os.Stat(path)
	if stat.Size() == 0 {
		_, err = db.Write([]byte("[]"))
		if err != nil {
			log.Fatalf("Problem writing to %s %v", path, err)
		}
		if db.Sync() != nil {
			log.Fatalf("Couldn't flush the file %s", path)
		}
	}
	return db, func() { db.Close() }
}

func NewFSPlayerStorage(db *os.File) *FileSystemPlayerStorage {
	storage := &FileSystemPlayerStorage{
		Db:     db,
		League: league.League{},
	}
	return storage
}

func (f *FileSystemPlayerStorage) PostPlayerScore(player string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	leag, err := f.getLeague()
	if err != nil {
		return err
	}

	if p := leag.Find(player); p != nil {
		p.Wins++
	} else {
		leag = append(leag, league.Player{Name: player, Wins: 1})
	}

	f.Db.Seek(0, io.SeekStart) // To always write from the beginning
	json.NewEncoder(f.Db).Encode(leag)
	return nil
}

func (f *FileSystemPlayerStorage) GetLeagueTable() (league.League, error) {
	league, err := f.getLeague()
	if err != nil {
		return nil, err
	}
	return league, nil
}

func (f *FileSystemPlayerStorage) GetPlayerScore(player string) (int, error) {
	score, err := f.getLeague()
	if err != nil {
		return 0, err
	}
	if p := score.Find(player); p != nil {
		return p.Wins, nil
	}
	return 0, errors.New(fmt.Sprintf("Requested player '%s' is missing", player))
}

func (f *FileSystemPlayerStorage) getLeague() (league.League, error) {
	f.Db.Seek(0, io.SeekStart) // To always read from the beginning
	var leag league.League
	err := json.NewDecoder(f.Db).Decode(&leag)

	if err != nil {
		err = fmt.Errorf("problem parsing league: %v\n", err)
	}
	return leag, err
}

func CreateTempFile(t testing.TB, initalData string) (io.ReadWriteSeeker, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	tmpfile.Write([]byte(initalData))
	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
