package fss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
)

type FileSystemPlayerStorage struct {
	mu     sync.Mutex
	Db     *json.Encoder
	League leaguedb.League
}

func FileSystemStorageFromFile(path string) (*FileSystemPlayerStorage, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("Problem opening %s %v", path, err)
	}

	close := func() {
		db.Close()
	}

	store, err := NewFSPlayerStorage(db)
	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}
	return store, close, nil
}

func initPlayersDBFile(db *os.File) error {
	stat, err := db.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", db.Name(), err)
	}

	if stat.Size() == 0 {
		_, err = db.Write([]byte("[]"))
		if err != nil {
			log.Fatalf("Problem writing to %s %v", db.Name(), err)
		}
		db.Seek(0, io.SeekStart)
	}
	return nil
}

func NewFSPlayerStorage(db *os.File) (*FileSystemPlayerStorage, error) {
	db.Seek(0, io.SeekStart)

	err := initPlayersDBFile(db)
	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := leaguedb.NewLeague(db)
	if err != nil {
		return nil, fmt.Errorf("problem loading player storage from file %s, %v", db.Name(), err)
	}

	return &FileSystemPlayerStorage{
		Db:     json.NewEncoder(&tape{file: db}),
		League: league,
	}, nil
}

func (f *FileSystemPlayerStorage) PostPlayerScore(player string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if p := f.League.Find(player); p != nil {
		p.Wins++
	} else {
		f.League = append(f.League, leaguedb.Player{Name: player, Wins: 1})
	}
	// TODO: fix the issue related to deleting players (Though it's not implemented yet).
	// If file length will decrease compare to initial state it will break everything
	f.Db.Encode(f.League)
	return nil
}

func (f *FileSystemPlayerStorage) GetLeagueTable() (leaguedb.League, error) {
	sort.Slice(f.League, func(i, j int) bool {
		return f.League[i].Wins > f.League[j].Wins
	})
	return f.League, nil
}

func (f *FileSystemPlayerStorage) GetPlayerScore(player string) (int, error) {
	if p := f.League.Find(player); p != nil {
		return p.Wins, nil
	}
	return 0, errors.New(fmt.Sprintf("Requested player '%s' is missing", player))
}

func CreateTempFile(t testing.TB, initalData string) (*os.File, func()) {
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
