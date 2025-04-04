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
	BaseTime               = 5
	NumPlayerPrompt        = "Please enter the number of players: "
	WrongPlayerInputErrMsg = "Bad value received for number of players, please try again with a number\n"
)

type ScheduledAlert struct {
	At     time.Duration
	Amount int
}

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

type GameSpy struct {
	StartCalled       bool
	StartedCalledWith int
	BlindAlert        []byte

	FinishedCalled   bool
	FinishCalledWith string
}

func (g *GameSpy) Start(numberOfPlayers int, to io.Writer) {
	g.StartCalled = true
	g.StartedCalledWith = numberOfPlayers
	to.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
	g.FinishCalledWith = winner
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{at, amount})
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (c *CLI) PlayPoker() {
	fmt.Fprint(c.out, NumPlayerPrompt)

	trimmedPrompt := strings.Trim(c.readInput(), "\n")
	numOfPlayers, err := strconv.Atoi(trimmedPrompt)
	if err != nil {
		fmt.Fprint(c.out, WrongPlayerInputErrMsg)
		return
	}

	c.game.Start(numOfPlayers, c.out)

	userInput := c.readInput()
	c.game.Finish(getTheName(userInput))
}

func (c *CLI) readInput() string {
	c.in.Scan()
	return c.in.Text()
}

func (c *CLI) awaitNumOfPlayersPrompt() int {
	for {
		fmt.Fprint(c.out, NumPlayerPrompt)

		trimmedPrompt := strings.Trim(c.readInput(), "\n")
		num, err := strconv.Atoi(trimmedPrompt)
		if err != nil {
			fmt.Fprint(c.out, "Bad value received for number of players, please try again with a number\n")
			continue
		}
		return num
	}
}

func getTheName(input string) string {
	return strings.Replace(input, " wins", "", 1)
}
