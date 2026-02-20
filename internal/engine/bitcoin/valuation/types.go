package valuation

import "crypto-api/internal/models"

// State represents the computed valuation snapshot
// derived from BTC price and global M2 liquidity.
type State struct {
	Meta             models.Meta `json:"meta"`
	BtcPrice         float64     `json:"btcPrice"`
	M2SupplyBillions float64     `json:"m2SupplyBillions"`
	Ratio            float64     `json:"ratio"`
	Description      string      `json:"description"`
}
