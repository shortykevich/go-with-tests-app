package fss

import (
	"slices"
	"testing"

	"github.com/shortykevich/go-with-tests-app/db/league"
)

func TestFileSystemStorage(t *testing.T) {
	t.Run("league from a reader", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStorage{Db: db}

		got, err := store.GetLeagueTable()
		if err != nil {
			t.Fatal(err)
		}
		want := league.League{
			{Name: "Cleo", Wins: 10},
			{Name: "Chris", Wins: 33},
		}
		if !slices.Equal(got, want) {
			t.Errorf("players table is wrong, got %q, want %q", got, want)
		}

		got, err = store.GetLeagueTable()
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got, want) {
			t.Errorf("players table is wrong, got %q, want %q", got, want)
		}

	})

	t.Run("get player score", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStorage{Db: db}

		got, err := store.GetPlayerScore("Chris")
		if err != nil {
			t.Fatal(err)
		}
		want := 33

		AssertPlayerScore(t, got, want)
	})

	t.Run("store wins for existing players", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStorage{Db: db}

		err := store.PostPlayerScore("Chris")
		if err != nil {
			t.Fatal(err)
		}

		got, err := store.GetPlayerScore("Chris")
		if err != nil {
			t.Fatal(err)
		}
		want := 34

		AssertPlayerScore(t, got, want)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		db, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()

		store := FileSystemPlayerStorage{Db: db}

		store.PostPlayerScore("Pepper")

		got, err := store.GetPlayerScore("Pepper")
		if err != nil {
			t.Fatal(err)
		}
		want := 1
		AssertPlayerScore(t, got, want)
	})
}

func AssertPlayerScore(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
