package poker

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.amount, c.at), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				alert := blindAlerter.alerts[i]

				if alert.amount != c.amount {
					t.Errorf("got amount %d, want %d", alert.amount, c.amount)
				}
				if alert.at != c.at {
					t.Errorf("got scheduled time of %v, want %v", alert.at, c.at)
				}
			})
		}
	})
}
