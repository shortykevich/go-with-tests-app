package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
