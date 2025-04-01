package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/league"
	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
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

func (s *SpyStorage) GetLeagueTable() (league.League, error) {
	leag := make(league.League, 0, len(s.scores))
	for name, wins := range s.scores {
		leag = append(leag, league.Player{Name: name, Wins: wins})
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

func GetLeagueFromResponse(t testing.TB, body io.Reader) (leag league.League) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&leag); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func TestPlayersScores(t *testing.T) {
	storage := &SpyStorage{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		winCalls: []string{},
	}

	server := NewPlayersScoreServer(storage)

	tests := []struct {
		name               string
		player             string
		expectedHTTPStatus int
		expectedScore      string
	}{
		{
			name:               "Returns Pepper's score",
			player:             "Pepper",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "20",
		},
		{
			name:               "Returns Floyd's score",
			player:             "Floyd",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "10",
		},
		{
			name:               "Returns 404 on missing players",
			player:             "Apollo",
			expectedHTTPStatus: http.StatusNotFound,
			expectedScore:      "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := NewGetScoreRequest(tt.player)
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			tutils.AssertStatus(t, resp.Code, tt.expectedHTTPStatus)
			tutils.AssertResponseBody(t, resp.Body.String(), tt.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	storage := &SpyStorage{
		scores:   map[string]int{},
		winCalls: []string{},
	}
	server := NewPlayersScoreServer(storage)

	t.Run("Records wins when POST", func(t *testing.T) {
		player := "Pepper"

		req := NewPostRequest("Pepper")
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		tutils.AssertStatus(t, resp.Code, http.StatusAccepted)

		if len(storage.winCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(storage.winCalls), 1)
		}
		if storage.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", storage.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	storage := &SpyStorage{
		scores: map[string]int{
			"Alice": 15,
			"Bill":  10,
		},
	}
	server := NewPlayersScoreServer(storage)

	t.Run("get request on /league", func(t *testing.T) {
		req := NewLeagueRequest(t)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		want, _ := server.storage.GetLeagueTable()
		got := GetLeagueFromResponse(t, resp.Body)

		tutils.AssertStatus(t, resp.Code, http.StatusOK)
		tutils.AssertLeague(t, got, want)
		tutils.AssertContentType(t, *resp, jsonContentType)
	})
}
