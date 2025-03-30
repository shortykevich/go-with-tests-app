package main

import (
	"encoding/json"
	"log"
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
	http.Handler
}

func NewPlayersScoreServer(storage PlayersStorage) *PlayersScoreServer {
	serv := &PlayersScoreServer{
		Storage: storage,
	}

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(serv.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(serv.playersHandler))

	serv.Handler = router

	return serv
}

func (p *PlayersScoreServer) postWin(w http.ResponseWriter, name string) {
	if err := p.Storage.PostPlayerScore(name); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayersScoreServer) getScore(w http.ResponseWriter, name string) {
	v, err := p.Storage.GetPlayerScore(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write([]byte(strconv.Itoa(v)))
}

func (p *PlayersScoreServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	tt := []Player{
		{
			Name: "Bill",
			Wins: 10,
		},
		{
			Name: "Alice",
			Wins: 15,
		},
	}
	err := json.NewEncoder(w).Encode(tt)
	if err != nil {
		log.Printf("Unable to parse input: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (p *PlayersScoreServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.postWin(w, player)
	case http.MethodGet:
		p.getScore(w, player)
	}
}
