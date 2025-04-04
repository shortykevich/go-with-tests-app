package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	fss "github.com/shortykevich/go-with-tests-app/db/fs_storage"
	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	db, cleanDatabase := fss.CreateTempFile(t, "")
	defer cleanDatabase()
	store, err := fss.NewFSPlayerStorage(db)

	tutils.AssertNoError(t, err)

	server := NewPlayersScoreServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostRequest(player))

	t.Run("get score", func(t *testing.T) {

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))
		tutils.AssertStatus(t, response, http.StatusOK)

		tutils.AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest(t))
		tutils.AssertStatus(t, response, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := leaguedb.League{
			{Name: "Pepper", Wins: 3},
		}
		tutils.AssertLeague(t, got, want)
	})
}
