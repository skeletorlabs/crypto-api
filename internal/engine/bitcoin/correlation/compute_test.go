package correlation

import (
	"testing"
	"time"
)

func TestCompute_PerfectCorrelation(t *testing.T) {
	now := time.Now()
	// BTC e M2 going up together
	m2 := []DataPoint{
		{Date: now, Value: 100},
		{Date: now.AddDate(0, 0, -1), Value: 90},
	}
	btc := []DataPoint{
		{Date: now, Value: 50000},
		{Date: now.AddDate(0, 0, -1), Value: 45000},
	}

	result, err := Compute(m2, btc)
	if err != nil {
		t.Fatalf("Compute failed: %v", err)
	}

	// Correlation should be positive and high
	if result.Coefficient <= 0.9 {
		t.Errorf("expected high positive correlation, got %f", result.Coefficient)
	}
}

func TestCompute_EmptyData(t *testing.T) {
	t.Run("Empty slices", func(t *testing.T) {
		_, err := Compute([]DataPoint{}, []DataPoint{})
		if err == nil {
			t.Error("expected error for empty data sets, got nil")
		}
	})

	t.Run("Insufficient data (one point only)", func(t *testing.T) {
		now := time.Now()
		m2 := []DataPoint{{Date: now, Value: 100}}
		btc := []DataPoint{{Date: now, Value: 50000}}

		_, err := Compute(m2, btc)
		if err == nil {
			t.Error("expected error for insufficient data (need at least 2 points)")
		}
	})
}
