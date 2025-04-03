package poker

import (
	"fmt"
	"strings"
	"testing"

	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

var dummySpyAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {
	cases := []string{
		"Chris",
		"Cleo",
	}
	for _, name := range cases {
		t.Run(fmt.Sprintf("record %s win from user input", name), func(t *testing.T) {
			input := strings.NewReader(fmt.Sprintf("%s wins\n", name))
			storage := tutils.NewSpyStorage()
			cli := NewCLI(storage, input, dummySpyAlerter)

			cli.PlayPoker()
			tutils.AssertPlayerWin(t, storage, name)
		})
	}
}

func TestBlindAerter(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStorage := tutils.NewSpyStorage()
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStorage, in, blindAlerter)
		cli.PlayPoker()

		if len(blindAlerter.alerts) != 1 {
			t.Fatal("expected a blind alert to be scheduled")
		}
	})
}
