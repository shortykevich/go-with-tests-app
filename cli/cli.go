package poker

import (
	"bufio"
	"io"
	"strings"

	"github.com/shortykevich/go-with-tests-app/webserver"
)

type CLI struct {
	storage webserver.PlayersStorage
	in      *bufio.Scanner
}

func NewCLI(store webserver.PlayersStorage, in io.Reader) *CLI {
	return &CLI{
		storage: store,
		in:      bufio.NewScanner(in),
	}
}

func (c *CLI) PlayPoker() {
	userInput := c.readInput()
	c.storage.PostPlayerScore(getTheName(userInput))
}

func (c *CLI) readInput() string {
	c.in.Scan()
	return c.in.Text()
}

func getTheName(input string) string {
	return strings.Replace(input, " wins", "", 1)
}
