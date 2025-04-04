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
	StartCalled  bool
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartCalled = true
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
			storage := tutils.NewStubStorage()
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
		playerStorage := tutils.NewStubStorage()
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
		playerStorage := tutils.NewStubStorage()
		out := &bytes.Buffer{}
		in := userInput("7")
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
		playerStorage := tutils.NewStubStorage()
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
		in := userInput("7")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, stdout.String())
		assertGameStartedWith(t, game, 7)
	})

	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := userInput("Pies")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, numPlayerPrompt, wrongPlayerInputErrMsg)
		assertGameNotStarted(t, game)
	})

	t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
		game := &GameSpy{}
		stdout := &bytes.Buffer{}

		in := userInput("3", "Chris wins")
		cli := NewCLI(in, stdout, game)

		cli.PlayPoker()

		assertMessagesSentToUser(t, stdout, numPlayerPrompt)
		assertGameStartedWith(t, game, 3)
		assertFinishCalledWith(t, game, "Chris")
	})

	t.Run("start game with 8 players and record 'Cleo' as winner", func(t *testing.T) {
		game := &GameSpy{}

		in := userInput("8", "Cleo wins")
		cli := NewCLI(in, dummyStdOut, game)

		cli.PlayPoker()

		assertGameStartedWith(t, game, 8)
		assertFinishCalledWith(t, game, "Cleo")
	})
}

func TestGame_Finish(t *testing.T) {
	store := tutils.NewStubStorage()
	game := NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	tutils.AssertPlayerWin(t, store, winner)
}

func assertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func assertFinishCalledWith(t testing.TB, game *GameSpy, want string) {
	t.Helper()
	if game.FinishedWith != want {
		t.Errorf("got %s but expected %s", game.FinishedWith, want)
	}
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func assertGameStartedWith(t testing.TB, game *GameSpy, want int) {
	t.Helper()
	if game.StartedWith != want {
		t.Errorf("wanted Start called with %d but got %d", want, game.StartedWith)
	}
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

func userInput(msgs ...string) io.Reader {
	return strings.NewReader(strings.Join(msgs, "\n"))
}
