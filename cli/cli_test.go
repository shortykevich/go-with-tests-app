package poker

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
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
			cli := NewCLI(storage, input, dummyStdOut, dummySpyAlerter)

			cli.PlayPoker()
			tutils.AssertPlayerWin(t, storage, name)
		})
	}
}

func TestBlindAlerter(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStorage := tutils.NewSpyStorage()
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStorage, in, dummyStdOut, blindAlerter)
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

				assertScheduledAlert(t, alert, c)
			})
		}
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		playerStorage := tutils.NewSpyStorage()
		out := &bytes.Buffer{}
		cli := NewCLI(playerStorage, dummyStdIn, out, dummySpyAlerter)
		cli.PlayPoker()

		got := out.String()

		if got != numPlayerPrompt {
			t.Errorf("got %q, want %q", got, numPlayerPrompt)
		}
	})
}

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}
	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}
