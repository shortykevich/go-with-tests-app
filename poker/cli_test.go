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

func TestCLI(t *testing.T) {
	cases := []string{
		"Chris",
		"Cleo",
	}
	for _, name := range cases {
		t.Run(fmt.Sprintf("record %s win from user input", name), func(t *testing.T) {
			input := userInput("5", fmt.Sprintf("%s wins\n", name))
			storage := tutils.NewStubStorage()
			game := NewTexasHoldem(dummySpyAlerter, storage)
			cli := NewCLI(input, dummyStdOut, game)

			cli.PlayPoker()
			tutils.AssertPlayerWin(t, storage, name)
		})
	}
}

func TestBlindAlerter(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := userInput("5", "Chris wins")
		playerStorage := tutils.NewStubStorage()
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, playerStorage)
		cli := NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()

		cases := []ScheduledAlert{
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
		playerStorage := tutils.NewStubStorage()
		out := &bytes.Buffer{}
		in := userInput("7")
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, playerStorage)
		cli := NewCLI(in, out, game)

		cli.PlayPoker()

		got := out.String()
		want := NumPlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		playerStorage := tutils.NewStubStorage()
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, playerStorage)

		game.Start(7, io.Discard)

		cases := []ScheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		checkSchedulingCases(t, cases, blindAlerter)
	})

	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userInput("7")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		AssertMessagesSentToUser(t, stdout, stdout.String())
		AssertGameStartedWith(t, game, 7)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userInput("Pies")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		AssertMessagesSentToUser(t, stdout, NumPlayerPrompt, WrongPlayerInputErrMsg)
		AssertGameNotStarted(t, game)
	})

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userInput("3", "Chris wins")
		cli := NewCLI(in, stdout, game)

		cli.PlayPoker()

		AssertMessagesSentToUser(t, stdout, NumPlayerPrompt)
		AssertGameStartedWith(t, game, 3)
		AssertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userInput("8", "Cleo wins")
		cli := NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		AssertGameStartedWith(t, game, 8)
		AssertFinishCalledWith(t, game, "Cleo")
	})
}

func TestGame_Finish(t *testing.T) {
	store := tutils.NewStubStorage()
	game := NewTexasHoldem(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	tutils.AssertPlayerWin(t, store, winner)
}

func checkSchedulingCases(t *testing.T, cases []ScheduledAlert, alerter *SpyBlindAlerter) {
	t.Helper()
	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {
			if len(alerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, alerter.alerts)
			}
			got := alerter.alerts[i]
			AssertScheduledAlert(t, got, want)
		})
	}
}

func userInput(msgs ...string) io.Reader {
	return strings.NewReader(strings.Join(msgs, "\n"))
}
