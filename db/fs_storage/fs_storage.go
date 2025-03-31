package fss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/league"
)

type FileSystemPlayerStorage struct {
	mu sync.Mutex
	Db io.ReadWriteSeeker
}

func NewFSPlayerStorage(db *os.File) *FileSystemPlayerStorage {
	storage := &FileSystemPlayerStorage{Db: db}
	storage.Db.Write([]byte("[]"))
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
		err = fmt.Errorf("problem parsing league: %v", err)
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
