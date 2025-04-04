package tutils

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
)

type StubStorage struct {
	mu       sync.Mutex
	Scores   map[string]int
	WinCalls []string
}

func NewStubStorage() *StubStorage {
	return &StubStorage{
		Scores:   make(map[string]int),
		WinCalls: []string(nil),
	}
}

func (s *StubStorage) GetPlayerScore(name string) (int, error) {
	v, ok := s.Scores[name]
	if !ok {
		return 0, errors.New(fmt.Sprintf("There no player with '%s' name!\n", name))
	}
	return v, nil
}

func (s *StubStorage) PostPlayerScore(name string) error {
	s.RecordWin(name)
	return nil
}

func (s *StubStorage) RecordWin(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Scores[name]++
	s.WinCalls = append(s.WinCalls, name)
}

func (s *StubStorage) GetLeagueTable() (leaguedb.League, error) {
	leag := make(leaguedb.League, 0, len(s.Scores))
	for name, wins := range s.Scores {
		leag = append(leag, leaguedb.Player{Name: name, Wins: wins})
	}
	sort.Slice(leag, func(i, j int) bool {
		return leag[i].Wins > leag[j].Wins
	})
	return leag, nil
}
