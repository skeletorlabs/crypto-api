package models

import (
	"crypto-api/internal/trend"
)

type BitcoinNetworkResponse struct {
	Meta
	BlockHeight         int64              `json:"blockHeight"`
	HashrateTHs         float64            `json:"hashrateTHs"`
	Difficulty          float64            `json:"difficulty"`
	AvgBlockTimeSeconds float64            `json:"avgBlockTimeSeconds"`
	Trend               trend.NetworkTrend `json:"trend"`
}
