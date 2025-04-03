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
	baseTime               = 5
	numPlayerPrompt        = "Please enter the number of players: "
	wrongPlayerInputErrMsg = "Bad value received for number of players, please try again with a number\n"
)

type scheduledAlert struct {
	at     time.Duration
	amount int
}

type CLI struct {
	in   *bufio.Scanner
	out  io.Writer
	game Game
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

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:   bufio.NewScanner(in),
		out:  out,
		game: game,
	}
}

func (c *CLI) PlayPoker() {
	fmt.Fprint(c.out, numPlayerPrompt)

	trimmedPrompt := strings.Trim(c.readInput(), "\n")
	numOfPlayers, err := strconv.Atoi(trimmedPrompt)
	if err != nil {
		fmt.Fprint(c.out, wrongPlayerInputErrMsg)
		return
	}

	c.game.Start(numOfPlayers)

	userInput := c.readInput()
	c.game.Finish(getTheName(userInput))
}

func (c *CLI) readInput() string {
	c.in.Scan()
	return c.in.Text()
}

func (c *CLI) awaitNumOfPlayersPrompt() int {
	for {
		fmt.Fprint(c.out, numPlayerPrompt)

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
