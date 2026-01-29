package models

type BitcoinMempoolResponse struct {
	Count    int  `json:"count"`
	VSize    int  `json:"vsize"`
	TotalFee int  `json:"totalFee"`
	Cached   bool `json:"cached"`
}
