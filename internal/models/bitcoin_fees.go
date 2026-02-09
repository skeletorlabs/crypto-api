package models

type BitcoinFeesResponse struct {
	Meta   Meta `json:"meta"`
	Low    int  `json:"low"`
	Medium int  `json:"medium"`
	High   int  `json:"high"`
}
