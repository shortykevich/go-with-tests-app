package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/inmem"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := inmem.NewInMemoryStorage()
	server := NewPlayersScoreServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		AssertStatus(t, response.Code, http.StatusOK)

		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest(t))
		AssertStatus(t, response.Code, http.StatusOK)

		got := GetLeagueFromResponse(t, response.Body)
		want := []inmem.Player{
			{Name: "Pepper", Wins: 3},
		}
		AssertLeague(t, got, want)
	})
}
