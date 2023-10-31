package formatting_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/formatting"
)

func TestParseTime(t *testing.T) {
	now := time.Now()
	cc := []struct {
		s       string
		want    time.Time
		wantErr bool
	}{
		{
			s:    "01.02.2023",
			want: time.Date(2023, time.February, 1, 0, 0, 0, 0, now.Location()),
		},
		{
			s:    "10:30",
			want: time.Date(now.Year(), now.Month(), now.Day(), 10, 30, 0, 0, now.Location()),
		},
		{
			s:    "01.02.2023 10:30",
			want: time.Date(2023, time.February, 1, 10, 30, 0, 0, now.Location()),
		},

		{
			s:    "now",
			want: now,
		},
		{
			s:    "now-1d",
			want: now.Add(-time.Hour * 24),
		},
		{
			s:    "-1d",
			want: now.Add(-time.Hour * 24),
		},

		{
			s:       "foo",
			wantErr: true,
		},
		{
			s:       "now-foo",
			wantErr: true,
		},
		{
			s:       "01.01",
			wantErr: true,
		},
	}

	for _, c := range cc {
		t.Run(c.s, func(t *testing.T) {
			got, err := formatting.ParseTime(c.s, now)
			if c.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, c.want, got)
			}
		})
	}
}
