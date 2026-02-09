package models

type PriceResponse struct {
	Meta  Meta    `json:"meta"`
	Token string  `json:"token"`
	USD   float64 `json:"usd"`
}
