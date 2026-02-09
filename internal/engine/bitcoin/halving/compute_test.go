package halving

import (
	"testing"
)

func TestCompute_StartOfEpoch(t *testing.T) {
	// Block 840,000 is the start of a new Halving (April/2024)
	currentBlock := 840000
	avgTime := 10.0

	res := Compute(currentBlock, avgTime)

	// The next halving should be at 840,000 + 210,000
	if res.NextHalvingBlock != 1050000 {
		t.Errorf("expected next halving at 1050000, got %d", res.NextHalvingBlock)
	}

	// The progress should be 0% at the first block
	if res.ProgressPercent != 0 {
		t.Errorf("expected progress 0, got %f", res.ProgressPercent)
	}
}

func TestCompute_MidEpoch(t *testing.T) {
	// Block 945,000 is exactly halfway between 840k and 1050k
	currentBlock := 945000
	avgTime := 10.0

	res := Compute(currentBlock, avgTime)

	// The progress should be exactly 50%
	if res.ProgressPercent != 50.0 {
		t.Errorf("expected progress 50.0, got %f", res.ProgressPercent)
	}

	// Should have 105,000 blocks remaining
	if res.BlocksRemaining != 105000 {
		t.Errorf("expected 105000 blocks remaining, got %d", res.BlocksRemaining)
	}
}

func TestCompute_LastBlockOfEpoch(t *testing.T) {
	// Block 1.049.999 is the last block before the next halving
	currentBlock := 1049999
	avgTime := 10.0

	res := Compute(currentBlock, avgTime)

	if res.BlocksRemaining != 1 {
		t.Errorf("expected 1 block remaining, got %d", res.BlocksRemaining)
	}

	// The progress must be 100%
	if res.ProgressPercent < 99.9 {
		t.Errorf("expected progress near 100, got %f", res.ProgressPercent)
	}
}

func TestCompute(t *testing.T) {
	// Test for current era (block 840,000+)
	state := Compute(840000, 10.0)
	if state.CurrentSubsidy != 3.125 {
		t.Errorf("Expected 3.125, got %f", state.CurrentSubsidy)
	}
}
