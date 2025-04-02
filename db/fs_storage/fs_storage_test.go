package fss

import (
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

func TestFileSystemStorage(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFSPlayerStorage(db)
		tutils.AssertNoError(t, err)

		got, err := store.GetLeagueTable()
		tutils.AssertNoError(t, err)
		t.Logf("DEBUG: %v", got)
		want := leaguedb.League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}
		tutils.AssertLeague(t, got, want)

		got, err = store.GetLeagueTable()
		tutils.AssertNoError(t, err)
		tutils.AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFSPlayerStorage(db)
		tutils.AssertNoError(t, err)

		got, err := store.GetPlayerScore("Chris")
		tutils.AssertNoError(t, err)

		want := 33
		tutils.AssertPlayerScore(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFSPlayerStorage(db)
		tutils.AssertNoError(t, err)

		err = store.PostPlayerScore("Chris")
		tutils.AssertNoError(t, err)

		got, err := store.GetPlayerScore("Chris")
		if err != nil {
			t.Fatal(err)
		}
		want := 34

		tutils.AssertPlayerScore(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFSPlayerStorage(db)
		tutils.AssertNoError(t, err)

		store.PostPlayerScore("Pepper")

		got, err := store.GetPlayerScore("Pepper")
		tutils.AssertNoError(t, err)

		want := 1
		tutils.AssertPlayerScore(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFSPlayerStorage(database)

		tutils.AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store, err := NewFSPlayerStorage(db)
		tutils.AssertNoError(t, err)

		got, err := store.GetLeagueTable()
		tutils.AssertNoError(t, err)

		want := leaguedb.League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}
		tutils.AssertLeague(t, got, want)

		got, err = store.GetLeagueTable()
		tutils.AssertNoError(t, err)
		tutils.AssertLeague(t, got, want)
	})
}
