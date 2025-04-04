package poker

import (
	"bytes"
	"strings"
	"testing"
)

func AssertGameNotStarted(t testing.TB, game *GameSpy) {
	t.Helper()
	if game.StartCalled {
		t.Errorf("game should not have started")
	}
}

func AssertFinishCalledWith(t testing.TB, game *GameSpy, want string) {
	t.Helper()
	if game.FinishCalledWith != want {
		t.Errorf("got %s but expected %s", game.FinishCalledWith, want)
	}
}

func AssertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}

func AssertGameStartedWith(t testing.TB, game *GameSpy, want int) {
	t.Helper()
	if game.StartedCalledWith != want {
		t.Errorf("wanted Start called with %d but got %d", want, game.StartedCalledWith)
	}
}

func AssertScheduledAlert(t testing.TB, got, want ScheduledAlert) {
	if got.Amount != want.Amount {
		t.Errorf("got amount %d, want %d", got.Amount, want.Amount)
	}
	if got.At != want.At {
		t.Errorf("got scheduled time of %v, want %v", got.At, want.At)
	}
}
