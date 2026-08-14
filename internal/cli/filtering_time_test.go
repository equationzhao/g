package cli

import (
	"testing"
	"time"
)

func TestCombineTodayWithClock(t *testing.T) {
	parsed, err := time.ParseInLocation("15:04", "14:30", time.Local)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 12, 18, 45, 0, 0, time.Local)
	got := combineTodayWithClock(parsed, now)
	want := time.Date(2026, 8, 12, 14, 30, 0, 0, time.Local)
	if !got.Equal(want) {
		t.Fatalf("combineTodayWithClock = %v, want %v", got, want)
	}
}
