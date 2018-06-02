package report

import (
	"github.com/go-fed/activity/pub"
	"time"
)

var _ pub.Clock = &LocalClock{}

// Uses time.Now directly. Implementations may use clocks configurable to a
// specific time zone, etc.
type LocalClock struct{}

func (l *LocalClock) Now() time.Time {
	return time.Now()
}
