package main

import (
	"errors"
	"fmt"
	"sync"
)

type InMemStorage struct {
	mu       sync.Mutex
	Scores   map[string]int
	winCalls []string
}

func NewInMemoryStorage() *InMemStorage {
	return &InMemStorage{
		Scores:   map[string]int{},
		winCalls: []string{},
	}
}

func (ims *InMemStorage) GetPlayerScore(name string) (int, error) {
	v, ok := ims.Scores[name]
	if !ok {
		return 0, errors.New(fmt.Sprintf("Player with '%s' name not found\n", name))
	}
	return v, nil
}

func (ims *InMemStorage) PostPlayerScore(name string) error {
	ims.RecordWin(name)
	return nil
}

func (ims *InMemStorage) RecordWin(name string) {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.Scores[name]++
	ims.winCalls = append(ims.winCalls, name)
}
