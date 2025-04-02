package poker

import (
	"bufio"
	"strings"
	"testing"

	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

func TestCLI(t *testing.T) {
	t.Run("record chris win from user input", func(t *testing.T) {
		input := strings.NewReader("Chris wins\n")
		storage := &tutils.SpyStorage{
			Scores:   make(map[string]int),
			WinCalls: []string{},
		}
		cli := &CLI{
			storage: storage,
			in:      bufio.NewScanner(input),
		}
		cli.PlayPoker()
		tutils.AssertPlayerWin(t, storage, "Chris")
	})
	t.Run("record cleo win from user input", func(t *testing.T) {
		input := strings.NewReader("Cleo wins\n")
		storage := &tutils.SpyStorage{
			Scores:   make(map[string]int),
			WinCalls: []string{},
		}
		cli := &CLI{
			storage: storage,
			in:      bufio.NewScanner(input),
		}
		cli.PlayPoker()
		tutils.AssertPlayerWin(t, storage, "Cleo")
	})
}
