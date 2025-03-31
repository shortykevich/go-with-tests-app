package jsontuil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/shortykevich/go-with-tests-app/db/inmem"
)

func NewLeagueFromReader(r io.Reader) ([]inmem.Player, error) {
	var league []inmem.Player
	err := json.NewDecoder(r).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league: %v", err)
	}
	return league, err
}

func GetPlayerScoreFromReader(r io.Reader, name string) (int, error) {
	league, _ := NewLeagueFromReader(r)
	for _, v := range league {
		if v.Name == name {
			return v.Wins, nil
		}
	}
	return 0, errors.New("Requested player is missing")
}
