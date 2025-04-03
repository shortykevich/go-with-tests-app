package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shortykevich/go-with-tests-app/db/leaguedb"
)

type BlindAlerter interface {
	ScheduleAlertAt(time.Duration, int)
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

type CLI struct {
	storage leaguedb.PlayersStorage
	in      *bufio.Scanner
	alerter BlindAlerter
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

func NewCLI(store leaguedb.PlayersStorage, in io.Reader, alerter BlindAlerter) *CLI {
	return &CLI{
		storage: store,
		in:      bufio.NewScanner(in),
		alerter: alerter,
	}
}

func (c *CLI) PlayPoker() {
	c.scheduleBlindAlerts()

	userInput := c.readInput()
	c.storage.PostPlayerScore(getTheName(userInput))
}

func (c *CLI) scheduleBlindAlerts() {
	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second

	for _, blind := range blinds {
		c.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + 10*time.Minute
	}
}

func (c *CLI) readInput() string {
	c.in.Scan()
	return c.in.Text()
}

func getTheName(input string) string {
	return strings.Replace(input, " wins", "", 1)
}
