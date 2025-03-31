package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/inmem"
)

type SpyStorage struct {
	mu       sync.Mutex
	scores   map[string]int
	winCalls []string
}

func (s *SpyStorage) GetPlayerScore(name string) (int, error) {
	v, ok := s.scores[name]
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
	s.scores[name]++
	s.winCalls = append(s.winCalls, name)
}

func (s *SpyStorage) GetLeagueTable() ([]inmem.Player, error) {
	league := make([]inmem.Player, 0, len(s.scores))
	for name, wins := range s.scores {
		league = append(league, inmem.Player{Name: name, Wins: wins})
	}
	return league, nil
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

func GetLeagueFromResponse(t testing.TB, body io.Reader) (league []inmem.Player) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func AssertContentType(t testing.TB, response httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %v, got %v", want, response.Result().Header)
	}
}

func AssertLeague(t testing.TB, got, want []inmem.Player) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("players table is wrong, got %q, want %q", got, want)
	}
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
