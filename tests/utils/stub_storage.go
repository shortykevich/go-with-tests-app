package tutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
)

type SpyStorage struct {
	mu       sync.Mutex
	Scores   map[string]int
	WinCalls []string
}

func (s *SpyStorage) GetPlayerScore(name string) (int, error) {
	v, ok := s.Scores[name]
	if !ok {
		return 0, errors.New(fmt.Sprintf("There no player with '%s' name!\n", name))
	}
	return v, nil
}

func (s *SpyStorage) PostPlayerScore(name string) error {
	s.RecordWin(name)
	return nil
}

func (s *SpyStorage) RecordWin(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Scores[name]++
	s.WinCalls = append(s.WinCalls, name)
}

func (s *SpyStorage) GetLeagueTable() (leaguedb.League, error) {
	leag := make(leaguedb.League, 0, len(s.Scores))
	for name, wins := range s.Scores {
		leag = append(leag, leaguedb.Player{Name: name, Wins: wins})
	}
	return leag, nil
}

func NewPostRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func NewGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func NewLeagueRequest(t testing.TB) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/league", nil)
	if err != nil {
		t.Fatalf("Request failed with error: %v", err)
	}
	return req
}

func GetLeagueFromResponse(t testing.TB, body io.Reader) (leag leaguedb.League) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&leag); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}
