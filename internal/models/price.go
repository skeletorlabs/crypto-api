package models

type PriceResponse struct {
	Meta
	Token string  `json:"token"`
	USD   float64 `json:"usd"`
}
