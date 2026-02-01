package models

type BitcoinMempoolResponse struct {
	Meta
	Count    int `json:"count"`
	VSize    int `json:"vsize"`
	TotalFee int `json:"totalFee"`
}
