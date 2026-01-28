package models

type BitcoinFeesResponse struct {
	Low    int  `json:"low"`
	Medium int  `json:"medium"`
	High   int  `json:"high"`
	Cached bool `json:"cached"`
}
