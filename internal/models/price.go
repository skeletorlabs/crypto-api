package models

type PriceResponse struct {
	Meta
	Symbol string  `json:"symbol"`
	USD    float64 `json:"usd"`
}
