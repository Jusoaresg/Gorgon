package web

import (
	"testing"
	"time"
)

func TestMondayOfWeek(t *testing.T) {
	cases := []struct {
		name     string
		now      time.Time
		expected time.Time
	}{
		{
			name:     "middle of the week",
			now:      time.Date(2026, 7, 29, 15, 30, 0, 0, time.UTC),
			expected: time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "monday itself",
			now:      time.Date(2026, 7, 27, 8, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "sunday",
			now:      time.Date(2026, 8, 2, 10, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "month boundary",
			now:      time.Date(2026, 8, 1, 2, 0, 0, 0, time.UTC),
			expected: time.Date(2026, 7, 27, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "year boundary",
			now:      time.Date(2026, 1, 4, 2, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 12, 29, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := mondayOfWeek(tc.now)
			if !got.Equal(tc.expected) {
				t.Errorf("mondayOfWeek(%s) = %s, want %s", tc.now, got, tc.expected)
			}
		})
	}
}

func TestMondayOfWeekBuildsMondayToSunday(t *testing.T) {
	monday := mondayOfWeek(time.Date(2026, 8, 1, 2, 0, 0, 0, time.UTC))
	if got, want := monday.Weekday(), time.Monday; got != want {
		t.Fatalf("week start = %s, want %s", got, want)
	}
	for i := range 7 {
		day := monday.AddDate(0, 0, i)
		want := (time.Monday + time.Weekday(i)) % 7
		if day.Weekday() != want {
			t.Errorf("day %d = %s, want %s", i, day.Weekday(), want)
		}
	}
}

func TestComputeWeekStartWithParam(t *testing.T) {
	got := computeWeekStart("2026-08-01")
	want := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("computeWeekStart(\"2026-08-01\") = %s, want %s", got, want)
	}
}

func TestComputeWeekStartCurrentWeek(t *testing.T) {
	got := computeWeekStart("")
	want := mondayOfWeek(time.Now().UTC())
	if !got.Equal(want) {
		t.Errorf("computeWeekStart(\"\") = %s, want %s", got, want)
	}
}
