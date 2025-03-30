package main

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

func (s *SpyStorage) GetLeagueTable() ([]Player, error) {
	league := make([]Player, 0, len(s.scores))
	for name, wins := range s.scores {
		league = append(league, Player{Name: name, Wins: wins})
	}
	return league, nil
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
			req := newGetScoreRequest(tt.player)
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assertStatus(t, resp.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, resp.Body.String(), tt.expectedScore)
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

		req := newPostRequest("Pepper")
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		assertStatus(t, resp.Code, http.StatusAccepted)

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
		req := newLeagueRequest(t)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		want, _ := server.storage.GetLeagueTable()
		got := getLeagueFromResponse(t, resp.Body)

		assertStatus(t, resp.Code, http.StatusOK)
		assertLeague(t, got, want)
		assertContentType(t, *resp, jsonContentType)
	})
}

func newPostRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func newLeagueRequest(t testing.TB) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/league", nil)
	if err != nil {
		t.Fatalf("Request failed with error: %v", err)
	}
	return req
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&league); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func assertContentType(t testing.TB, response httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %v, got %v", want, response.Result().Header)
	}
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("players table is wrong, got %q, want %q", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
