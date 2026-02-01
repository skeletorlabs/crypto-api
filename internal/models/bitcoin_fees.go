package models

type BitcoinFeesResponse struct {
	Meta
	Low    int `json:"low"`
	Medium int `json:"medium"`
	High   int `json:"high"`
}
