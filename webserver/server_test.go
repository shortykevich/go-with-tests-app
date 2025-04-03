package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

func TestPlayersScores(t *testing.T) {
	storage := &tutils.SpyStorage{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		WinCalls: []string{},
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
			req := tutils.NewGetScoreRequest(tt.player)
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			tutils.AssertStatus(t, resp.Code, tt.expectedHTTPStatus)
			tutils.AssertResponseBody(t, resp.Body.String(), tt.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	storage := &tutils.SpyStorage{
		Scores:   map[string]int{},
		WinCalls: []string{},
	}
	server := NewPlayersScoreServer(storage)

	t.Run("Records wins when POST", func(t *testing.T) {
		player := "Pepper"

		req := tutils.NewPostRequest("Pepper")
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		tutils.AssertStatus(t, resp.Code, http.StatusAccepted)

		if len(storage.WinCalls) != 1 {
			t.Fatalf("got %d calls to RecordWin want %d", len(storage.WinCalls), 1)
		}
		if storage.WinCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", storage.WinCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	t.Run("get request on /league", func(t *testing.T) {
		storage := &tutils.SpyStorage{
			Scores: map[string]int{
				"Alice": 15,
				"Bill":  10,
			},
		}
		server := NewPlayersScoreServer(storage)

		req := tutils.NewLeagueRequest(t)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		want, err := server.storage.GetLeagueTable()
		tutils.AssertNoError(t, err)

		got := tutils.GetLeagueFromResponse(t, resp.Body)

		tutils.AssertStatus(t, resp.Code, http.StatusOK)
		tutils.AssertLeague(t, got, want)
		tutils.AssertContentType(t, *resp, jsonContentType)
	})
}
