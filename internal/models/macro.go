package models

import "time"

type MacroResponse struct {
	Meta     Meta      `json:"meta"`
	M2Supply M2Details `json:"m2_supply"`
}

type M2Details struct {
	Value    float64   `json:"value"`
	Unit     string    `json:"unit"`
	DateTime time.Time `json:"date_time"`
}
