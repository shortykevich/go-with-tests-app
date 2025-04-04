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
	"github.com/shortykevich/go-with-tests-app/poker"
	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

var (
	dummyGame = &poker.GameSpy{}
	tenMS     = 10 * time.Millisecond
)

func TestPlayersScores(t *testing.T) {
	storage := &tutils.StubStorage{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		WinCalls: []string{},
	}

	server, err := NewPlayersScoreServer(storage, dummyGame)
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
	server, err := NewPlayersScoreServer(storage, dummyGame)
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
		server, err := NewPlayersScoreServer(storage, dummyGame)
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
		server, err := NewPlayersScoreServer(tutils.NewStubStorage(), dummyGame)
		tutils.AssertNoError(t, err)

		req := newGameRequest()
		resp := httptest.NewRecorder()

		server.ServeHTTP(resp, req)

		tutils.AssertStatus(t, resp, http.StatusOK)
	})

	t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		storage := tutils.NewStubStorage()
		game := &poker.GameSpy{BlindAlert: []byte(wantedBlindAlert)}
		winner := "Ruth"

		server := httptest.NewServer(mustMakePlayerServer(t, storage, game))
		wsURL := fmt.Sprintf("ws%s/ws", strings.TrimPrefix(server.URL, "http"))
		ws := mustDialWS(t, wsURL)

		defer func() {
			ws.Close()
			server.Close()
		}()

		writeWSMessage(t, ws, "3")
		writeWSMessage(t, ws, winner)

		time.Sleep(tenMS)

		poker.AssertGameStartedWith(t, game, 3)
		poker.AssertFinishCalledWith(t, game, winner)
		within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
	})
}

func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
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

func mustMakePlayerServer(t *testing.T, store leaguedb.PlayersStorage, game poker.Game) *PlayersScoreServer {
	server, err := NewPlayersScoreServer(store, game)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
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

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}
