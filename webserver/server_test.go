package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

func TestPlayersScores(t *testing.T) {
	storage := &tutils.StubStorage{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		WinCalls: []string{},
	}

	server, err := NewPlayersScoreServer(storage)
	tutils.AssertNoError(t, err)

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

			tutils.AssertStatus(t, resp, tt.expectedHTTPStatus)
			tutils.AssertResponseBody(t, resp.Body.String(), tt.expectedScore)
		})
	}
}

func TestStoreWins(t *testing.T) {
	storage := &tutils.StubStorage{
		Scores:   map[string]int{},
		WinCalls: []string{},
	}
	server, err := NewPlayersScoreServer(storage)
	tutils.AssertNoError(t, err)

	t.Run("Records wins when POST", func(t *testing.T) {
		player := "Pepper"

		req := newPostRequest("Pepper")
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		tutils.AssertStatus(t, resp, http.StatusAccepted)

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
		storage := &tutils.StubStorage{
			Scores: map[string]int{
				"Alice": 15,
				"Bill":  10,
			},
		}
		server, err := NewPlayersScoreServer(storage)
		tutils.AssertNoError(t, err)

		req := newLeagueRequest(t)
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		want, err := server.storage.GetLeagueTable()
		tutils.AssertNoError(t, err)

		got := getLeagueFromResponse(t, resp.Body)

		tutils.AssertStatus(t, resp, http.StatusOK)
		tutils.AssertLeague(t, got, want)
		tutils.AssertContentType(t, *resp, jsonContentType)
	})

	t.Run("GET /game return 200", func(t *testing.T) {
		server, err := NewPlayersScoreServer(tutils.NewStubStorage())
		tutils.AssertNoError(t, err)

		req := newGameRequest()
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		tutils.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		storage := tutils.NewStubStorage()
		winner := "Ruth"
		handler, err := NewPlayersScoreServer(storage)
		tutils.AssertNoError(t, err)

		server := httptest.NewServer(handler)
		defer server.Close()

		wsURL := fmt.Sprintf("ws%s/ws", strings.TrimPrefix(server.URL, "http"))

		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, winner)
		time.Sleep(10 * time.Millisecond)
		tutils.AssertPlayerWin(t, storage, winner)
	})
}

func newGameRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/game", nil)
	return req
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

func getLeagueFromResponse(t testing.TB, body io.Reader) (leag leaguedb.League) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(&leag); err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func mustDialWS(t *testing.T, wsURL string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}
	return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, msg string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
