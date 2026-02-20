package halving

import (
	"testing"
	"time"
)

var fixedTime = time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)

func TestCompute_StartOfEpoch(t *testing.T) {
	currentBlock := 840000
	avgTime := 10.0

	res := Compute(currentBlock, avgTime, fixedTime)

	if res.NextHalvingBlock != 1050000 {
		t.Errorf("expected next halving at 1050000, got %d", res.NextHalvingBlock)
	}

	if res.ProgressPercent != 0 {
		t.Errorf("expected progress 0, got %f", res.ProgressPercent)
	}
}

func TestCompute_MidEpoch(t *testing.T) {
	currentBlock := 945000
	avgTime := 10.0

	res := Compute(currentBlock, avgTime, fixedTime)

	if res.ProgressPercent != 50.0 {
		t.Errorf("expected progress 50.0, got %f", res.ProgressPercent)
	}

	if res.BlocksRemaining != 105000 {
		t.Errorf("expected 105000 blocks remaining, got %d", res.BlocksRemaining)
	}
}

func TestCompute_LastBlockOfEpoch(t *testing.T) {
	currentBlock := 1049999
	avgTime := 10.0

	res := Compute(currentBlock, avgTime, fixedTime)

	if res.BlocksRemaining != 1 {
		t.Errorf("expected 1 block remaining, got %d", res.BlocksRemaining)
	}

	if res.ProgressPercent < 99.9 {
		t.Errorf("expected progress near 100, got %f", res.ProgressPercent)
	}
}

func TestCompute_CurrentSubsidy(t *testing.T) {
	state := Compute(840000, 10.0, fixedTime)

	if state.CurrentSubsidy != 3.125 {
		t.Errorf("Expected 3.125, got %f", state.CurrentSubsidy)
	}
}
