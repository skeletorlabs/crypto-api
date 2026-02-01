package models

type ChainsResponse struct {
	Meta
	Data []ChainResponse `json:"data"`
}

type ChainResponse struct {
	Name   string  `json:"name"`
	TVL    float64 `json:"tvl"`
	Symbol string  `json:"symbol"`
}
