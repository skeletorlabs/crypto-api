package models

type ChainResponse struct {
	Name   string  `json:"name"`
	TVL    float64 `json:"tvl"`
	Symbol string  `json:"symbol"`
}
