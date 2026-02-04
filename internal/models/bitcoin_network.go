package models

import (
	"crypto-api/internal/engine/bitcoin/halving"
	"crypto-api/internal/engine/bitcoin/trend"
)

type BitcoinNetworkResponse struct {
	Meta
	BlockHeight         int64         `json:"blockHeight"`
	HashrateTHs         float64       `json:"hashrateTHs"`
	Difficulty          float64       `json:"difficulty"`
	AvgBlockTimeSeconds float64       `json:"avgBlockTimeSeconds"`
	Trend               trend.Status  `json:"trend"`
	Halving             halving.State `json:"halving"`
}
