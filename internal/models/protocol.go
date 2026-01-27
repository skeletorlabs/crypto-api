package models

type ProtocolResponse struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	TVL      float64 `json:"tvl"`
	Chain    string  `json:"chain"`
	Category string  `json:"category"`
}
