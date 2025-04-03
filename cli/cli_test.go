package poker

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	tutils "github.com/shortykevich/go-with-tests-app/tests/utils"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}
var dummySpyAlerter = &SpyBlindAlerter{}

type GameSpy struct {
	StartedWith  int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func TestCLI(t *testing.T) {
	cases := []string{
		"Chris",
		"Cleo",
	}
	for _, name := range cases {
		t.Run(fmt.Sprintf("record %s win from user input", name), func(t *testing.T) {
			input := userInput("5", fmt.Sprintf("%s wins\n", name))
			storage := tutils.NewSpyStorage()
			game := NewGame(dummySpyAlerter, storage)
			cli := NewCLI(input, dummyStdOut, game)

			cli.PlayPoker()
			tutils.AssertPlayerWin(t, storage, name)
		})
	}
}

func TestBlindAlerter(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := userInput("5", "Chris wins")
		playerStorage := tutils.NewSpyStorage()
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, playerStorage)
		cli := NewCLI(in, dummyStdOut, game)
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

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		playerStorage := tutils.NewSpyStorage()
		out := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, playerStorage)
		cli := NewCLI(in, out, game)

		cli.PlayPoker()

		got := out.String()
		want := numPlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		playerStorage := tutils.NewSpyStorage()
		blindAlerter := &SpyBlindAlerter{}
		game := NewGame(blindAlerter, playerStorage)

		game.Start(7)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		gotPrompt := stdout.String()
		wantPrompt := numPlayerPrompt

		if gotPrompt != wantPrompt {
			t.Errorf("got %q, want %q", gotPrompt, wantPrompt)
		}

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}
	})
}

func TestGame_Finish(t *testing.T) {
	store := tutils.NewSpyStorage()
	game := NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	tutils.AssertPlayerWin(t, store, winner)
}

func userInput(msgs ...string) io.Reader {
	return strings.NewReader(strings.Join(msgs, "\n"))
}

func assertScheduledAlert(t testing.TB, got, want scheduledAlert) {
	if got.amount != want.amount {
		t.Errorf("got amount %d, want %d", got.amount, want.amount)
	}
	if got.at != want.at {
		t.Errorf("got scheduled time of %v, want %v", got.at, want.at)
	}
}

func checkSchedulingCases(t *testing.T, cases []scheduledAlert, alerter *SpyBlindAlerter) {
	t.Helper()
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(alerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, alerter.alerts)
			}
			got := alerter.alerts[i]
			assertScheduledAlert(t, got, want)
		})
	}
}
