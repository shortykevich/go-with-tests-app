package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
	"github.com/shortykevich/go-with-tests-app/db/league"
	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	db, cleanDatabase := fss.CreateTempFile(t, "")
	defer cleanDatabase()
	store, err := fss.NewFSPlayerStorage(db)

	tutils.AssertNoError(t, err)

	server := NewPlayersScoreServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostRequest(player))

	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		tutils.AssertStatus(t, response.Code, http.StatusOK)

		tutils.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest(t))
		tutils.AssertStatus(t, response.Code, http.StatusOK)

		got := GetLeagueFromResponse(t, response.Body)
		want := []league.Player{
			{Name: "Pepper", Wins: 3},
		}
		tutils.AssertLeague(t, got, want)
	})
}
