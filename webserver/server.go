package webserver

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
	"github.com/shortykevich/go-with-tests-app/poker"
)

const (
	jsonContentType  = "application/json"
	htmlTemplatePath = "game.html"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PlayersScoreServer struct {
	storage leaguedb.PlayersStorage
	http.Handler
	template *template.Template
	game     poker.Game
}

type playerServerWS struct {
	*websocket.Conn
}

func NewPlayersScoreServer(storage leaguedb.PlayersStorage, game poker.Game) (*PlayersScoreServer, error) {
	serv := &PlayersScoreServer{}

	tmpl, err := template.ParseFiles(htmlTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	serv.template = tmpl
	serv.storage = storage
	serv.game = game

	router := http.NewServeMux()
	router.Handle("/ws", http.HandlerFunc(serv.webSocket))
	router.Handle("/game", http.HandlerFunc(serv.newGameHandler))
	router.Handle("/league", http.HandlerFunc(serv.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(serv.playersHandler))

	serv.Handler = router

	return serv, nil
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
		log.Printf("Couldn't get Players table. Error occurred. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(players)
	if err != nil {
		log.Printf("Unable to parse Players table. Error occurred. %v", err)
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

func (p *PlayersScoreServer) newGameHandler(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}

func (p *PlayersScoreServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numOfPlayersPrompt := ws.WaitForMsg()
	numOfPlayers, _ := strconv.Atoi(numOfPlayersPrompt)
	p.game.Start(numOfPlayers, ws)

	winner := ws.WaitForMsg()
	p.game.Finish(string(winner))
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{Conn: conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}

func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
