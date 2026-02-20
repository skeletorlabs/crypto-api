package models

import "crypto-api/internal/engine/bitcoin/halving"

type BitcoinNetworkResponse struct {
	Meta                Meta          `json:"meta"`
	BlockHeight         int64         `json:"blockHeight"`
	HashrateTHs         float64       `json:"hashrateTHs"`
	Difficulty          float64       `json:"difficulty"`
	AvgBlockTimeSeconds float64       `json:"avgBlockTimeSeconds"`
	Trend               Status        `json:"trend"`
	Halving             halving.State `json:"halving"`
}
