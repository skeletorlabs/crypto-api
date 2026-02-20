package valuation

// Compute executes the valuation logic and returns the current State.
// Ratio represents BTC price relative to global M2 liquidity (in billions).
func Compute(btcPrice float64, m2SupplyBillions float64) State {
	ratio := 0.0

	if m2SupplyBillions > 0 {
		// Relationship between Bitcoin price and global liquidity
		ratio = btcPrice / m2SupplyBillions
	}

	return State{
		BtcPrice:         btcPrice,
		M2SupplyBillions: m2SupplyBillions,
		Ratio:            ratio,
		Description:      "Bitcoin price divided by M2 Money Supply (Billions)",
	}
}
