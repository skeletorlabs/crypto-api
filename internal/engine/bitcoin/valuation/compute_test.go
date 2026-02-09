package valuation

import (
	"testing"
)

func TestCompute_Logic(t *testing.T) {
	// Controlled scenario
	btcPrice := 60000.0
	m2Billions := 21000.0 // Eg: 21T dollars

	state := Compute(btcPrice, m2Billions)

	// Validate if the ratio is correctly calculated
	// If Ratio = btcPrice / m2Billions
	expectedRatio := btcPrice / m2Billions
	if state.Ratio != expectedRatio {
		t.Errorf("expected ratio %f, got %f", expectedRatio, state.Ratio)
	}

	if state.Description == "" {
		t.Error("description should not be empty")
	}

	// The BtcPrice should be correctly reflected in the state
	if state.BtcPrice != btcPrice {
		t.Errorf("expected price %f, got %f", btcPrice, state.BtcPrice)
	}
}
