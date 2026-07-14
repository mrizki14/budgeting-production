package report

import (
	"testing"
	"time"
)

func TestNormalizePeriodFallsBackToCurrentDate(t *testing.T) {
	now := time.Date(2026, time.July, 7, 0, 0, 0, 0, time.UTC)
	month, year := normalizePeriod(0, 2035, now)

	if month != 7 || year != 2026 {
		t.Fatalf("expected 7/2026, got %d/%d", month, year)
	}
}
