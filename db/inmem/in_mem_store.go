package inmem

import (
	"errors"
	"fmt"
	"sync"
)

type InMemStorage struct {
	mu       sync.Mutex
	scores   map[string]int
	winCalls []string
}

type Player struct {
	Name string
	Wins int
}

func NewInMemoryStorage() *InMemStorage {
	return &InMemStorage{
		scores:   map[string]int{},
		winCalls: []string{},
	}
}

func (ims *InMemStorage) GetPlayerScore(name string) (int, error) {
	v, ok := ims.scores[name]
	if !ok {
		return 0, errors.New(fmt.Sprintf("Player with '%s' name not found\n", name))
	}
	return v, nil
}

func (ims *InMemStorage) PostPlayerScore(name string) error {
	ims.RecordWin(name)
	return nil
}

func (ims *InMemStorage) GetLeagueTable() ([]Player, error) {
	league := make([]Player, 0, len(ims.scores))
	for name, wins := range ims.scores {
		league = append(league, Player{Name: name, Wins: wins})
	}
	return league, nil
}

func (ims *InMemStorage) RecordWin(name string) {
	ims.mu.Lock()
	defer ims.mu.Unlock()

	ims.scores[name]++
	ims.winCalls = append(ims.winCalls, name)
}
