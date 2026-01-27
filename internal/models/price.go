package models

type PriceResponse struct {
	Symbol string  `json:"symbol"`
	USD    float64 `json:"usd"`
	Cached bool    `json:"cached"`
}
