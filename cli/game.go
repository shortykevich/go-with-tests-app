package poker

import (
	"time"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
)

type Game struct {
	alerter BlindAlerter
	storage leaguedb.PlayersStorage
}

func NewGame(alerter BlindAlerter, storage leaguedb.PlayersStorage) *Game {
	return &Game{
		alerter: alerter,
		storage: storage,
	}
}

func (g *Game) Start(numOfPlayers int) {
	blindInc := time.Duration(baseTime+numOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime += blindInc
	}
}

func (g *Game) Finish(winner string) {
	g.storage.PostPlayerScore(winner)
}
