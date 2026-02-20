package models

import "time"

// IntelligencePrice represents a cached market price.
type IntelligencePrice struct {
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

// IntelligenceSnapshot represents a persisted structural intelligence state.
type IntelligenceSnapshot struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`

	// Market & Liquidity
	PriceUSD         float64 `json:"priceUsd"`
	M2SupplyBillions float64 `json:"m2SupplyBillions"`

	// Structural Metrics
	BTCM2Ratio  float64 `json:"btcM2Ratio"`
	Correlation float64 `json:"correlation"`

	// Network Metrics
	BlockHeight        int64   `json:"blockHeight"`
	HashrateTHs        float64 `json:"hashrateTHs"`
	Difficulty         float64 `json:"difficulty"`
	AvgBlockTime       float64 `json:"avgBlockTime"`
	NetworkHealthScore int     `json:"networkHealthScore"`
	TrendStatus        Status  `json:"trendStatus"`

	// Metadata
	SourceAttribution string `json:"sourceAttribution"`
}

type Status string

const (
	TrendImproving Status = "Improving"
	TrendStable    Status = "Stable"
	TrendWorsening Status = "Worsening"
)
