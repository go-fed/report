package report

import (
	"github.com/go-fed/activity/pub"
	"time"
)

var _ pub.Clock = &localClock{}

// Uses time.Now directly. Implementations may use clocks configurable to a
// specific time zone, etc.
type localClock struct{}

func (l *localClock) Now() time.Time {
	return time.Now()
}
