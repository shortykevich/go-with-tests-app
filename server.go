package main

import (
	"net/http"
	"strconv"
	"strings"
)

type PlayersStorage interface {
	GetPlayerScore(string) (int, error)
	PostPlayerScore(string) error
}

type PlayersScoreServer struct {
	Storage PlayersStorage
}

func (p *PlayersScoreServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		p.addWin(w, r)
	case http.MethodGet:
		p.getScore(w, r)
	}
}

func (p *PlayersScoreServer) addWin(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	if err := p.Storage.PostPlayerScore(player); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayersScoreServer) getScore(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	v, err := p.Storage.GetPlayerScore(player)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write([]byte(strconv.Itoa(v)))
}
