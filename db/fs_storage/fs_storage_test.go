package fss

import (
	"strings"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/inmem"
	"github.com/shortykevich/go-with-tests-app/webserver"
)

func TestFileSystemStorage(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		db := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStorage{db: db}

		got, err := store.GetLeagueTable()
		if err != nil {
			t.Fatal(err)
		}
		want := []inmem.Player{
			{Name: "Cleo", Wins: 10},
			{Name: "Chris", Wins: 33},
		}
		webserver.AssertLeague(t, got, want)

		got, err = store.GetLeagueTable()
		if err != nil {
			t.Fatal(err)
		}
		webserver.AssertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {
		db := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStorage{db: db}

		got, _ := store.GetPlayerScore("Chris")
		want := 33

		AssertPlayerScore(t, got, want)
	})
}

func AssertPlayerScore(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
