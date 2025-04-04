package poker

import (
	"fmt"
	"io"
	"time"
)

type BlindAlerter interface {
	ScheduleAlertAt(time.Duration, int, io.Writer)
}

type BlindAlerterFunc func(time.Duration, int, io.Writer)

func (b BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	b(duration, amount, to)
}

func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
