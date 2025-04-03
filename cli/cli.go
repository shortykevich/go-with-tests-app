package poker

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	baseTime        = 5
	numPlayerPrompt = "Please enter the number of players: "
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game *Game
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

func NewCLI(in io.Reader, out io.Writer, game *Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (c *CLI) PlayPoker() {
	fmt.Fprint(c.out, numPlayerPrompt)

	trimmedPrompt := strings.Trim(c.readInput(), "\n")
	numOfPlayers, _ := strconv.Atoi(trimmedPrompt)

	c.game.Start(numOfPlayers)

	userInput := c.readInput()
	c.game.Finish(getTheName(userInput))
}

func (c *CLI) readInput() string {
	c.in.Scan()
	return c.in.Text()
}

func getTheName(input string) string {
	return strings.Replace(input, " wins", "", 1)
}
