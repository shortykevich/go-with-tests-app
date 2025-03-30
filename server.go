package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const jsonContentType = "application/json"

type PlayersStorage interface {
	GetPlayerScore(string) (int, error)
	PostPlayerScore(string) error
	GetLeagueTable() ([]Player, error)
}

type PlayersScoreServer struct {
	storage PlayersStorage
	http.Handler
}

func NewPlayersScoreServer(storage PlayersStorage) *PlayersScoreServer {
	serv := &PlayersScoreServer{
		storage: storage,
	}

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(serv.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(serv.playersHandler))

	serv.Handler = router

	return serv
}

func (p *PlayersScoreServer) postWin(w http.ResponseWriter, name string) {
	if err := p.storage.PostPlayerScore(name); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayersScoreServer) getScore(w http.ResponseWriter, name string) {
	v, err := p.storage.GetPlayerScore(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write([]byte(strconv.Itoa(v)))
}

func (p *PlayersScoreServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	players, err := p.storage.GetLeagueTable()
	if err != nil {
		log.Printf("Couldn't get Players table. Error occurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(players)
	if err != nil {
		log.Printf("Unable to parse Players table. Error occurred: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
