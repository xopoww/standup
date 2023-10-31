package formatting

import (
	"errors"
	"strings"
	"time"

	"github.com/trstringer/go-systemd-time/pkg/systemdtime"
)

var ErrInvalidTimeFormat = errors.New("invalid time format")

func ParseTime(s string, now time.Time) (time.Time, error) {
	if s == "now" {
		return now, nil
	}

	// date + time
	if t, err := time.ParseInLocation("02.01.2006 15:04:05", s, now.Location()); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("02.01.2006 15:04", s, now.Location()); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("02.01.2006", s, now.Location()); err == nil {
		return t, nil
	}

	// time (assume today)
	if t, err := time.ParseInLocation("15:04:05", s, now.Location()); err == nil {
		return time.Date(
			now.Year(), now.Month(), now.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
			now.Location(),
		), nil
	}
	if t, err := time.ParseInLocation("15:04", s, now.Location()); err == nil {
		return time.Date(
			now.Year(), now.Month(), now.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
			now.Location(),
		), nil
	}

	// systemd adjustment (from now)
	if t, err := systemdtime.AdjustTime(now, strings.TrimPrefix(s, "now")); err == nil && !t.After(now) {
		return t, nil
	}

	return time.Time{}, ErrInvalidTimeFormat
}
