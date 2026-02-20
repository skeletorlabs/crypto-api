package models

import "time"

// PriceHistory represents a historical price record for an asset.
type PriceHistory struct {
	Timestamp time.Time `json:"timestamp"`
	Asset     string    `json:"asset"`
	PriceUSD  float64   `json:"price_usd"`
	Source    string    `json:"source"`
}
