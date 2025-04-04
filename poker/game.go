package poker

import (
	"io"
	"time"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
)

type Game interface {
	Start(int, io.Writer)
	Finish(string)
}

type TexasHoldem struct {
	alerter BlindAlerter
	storage leaguedb.PlayersStorage
}

func NewTexasHoldem(alerter BlindAlerter, storage leaguedb.PlayersStorage) *TexasHoldem {
	return &TexasHoldem{
		alerter: alerter,
		storage: storage,
	}
}

func (g *TexasHoldem) Start(numOfPlayers int, to io.Writer) {
	blindInc := time.Duration(BaseTime+numOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		g.alerter.ScheduleAlertAt(blindTime, blind, to)
		blindTime += blindInc
	}
}

func (g *TexasHoldem) Finish(winner string) {
	g.storage.PostPlayerScore(winner)
}
