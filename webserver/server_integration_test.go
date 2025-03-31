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

	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)

		assertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest(t))
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []inmem.Player{
			{Name: "Pepper", Wins: 3},
		}
		assertLeague(t, got, want)
	})
}
