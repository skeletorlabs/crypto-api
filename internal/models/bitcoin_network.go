package models

type BitcoinNetworkResponse struct {
	BlockHeight         int64   `json:"blockHeight"`
	HashrateTHs         float64 `json:"hashrateTHs"`
	Difficulty          float64 `json:"difficulty"`
	AvgBlockTimeSeconds float64 `json:"avgBlockTimeSeconds"`
	Cached              bool    `json:"cached"`
}
