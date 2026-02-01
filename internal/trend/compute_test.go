package trend

import (
	"testing"
	"time"
)

func TestComputeTrend_Improving(t *testing.T) {
	base := time.Now()

	snaps := []Snapshot{
		{
			Timestamp:       base,
			AvgBlockTimeSec: 720, // 12 min
		},
		{
			Timestamp:       base.Add(10 * time.Minute),
			AvgBlockTimeSec: 600, // 10 min (improved)
		},
	}

	trend := ComputeTrend(snaps)

	if trend != TrendImproving {
		t.Fatalf("expected %s, got %s", TrendImproving, trend)
	}
}

func TestComputeTrend_Stable(t *testing.T) {
	base := time.Now()

	snaps := []Snapshot{
		{
			Timestamp:       base,
			AvgBlockTimeSec: 600,
		},
		{
			Timestamp:       base.Add(10 * time.Minute),
			AvgBlockTimeSec: 610, // small variation
		},
	}

	trend := ComputeTrend(snaps)

	if trend != TrendStable {
		t.Fatalf("expected %s, got %s", TrendStable, trend)
	}
}

func TestComputeTrend_Worsening(t *testing.T) {
	base := time.Now()

	snaps := []Snapshot{
		{
			Timestamp:       base,
			AvgBlockTimeSec: 600, // 10 min
		},
		{
			Timestamp:       base.Add(10 * time.Minute),
			AvgBlockTimeSec: 900, // 15 min (it pretty damn worse)
		},
	}

	trend := ComputeTrend(snaps)

	if trend != TrendWorsening {
		t.Fatalf("expected %s, got %s", TrendWorsening, trend)
	}
}
