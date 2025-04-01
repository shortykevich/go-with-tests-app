package league

import (
	"encoding/json"
	"fmt"
	"io"
)

type Player struct {
	Name string
	Wins int
}

type League []Player

func NewLeague(r io.Reader) (League, error) {
	var leag League
	err := json.NewDecoder(r).Decode(&leag)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}
	return leag, err
}

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
