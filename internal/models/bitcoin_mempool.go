package models

type BitcoinMempoolResponse struct {
	Meta     Meta `json:"meta"`
	Count    int  `json:"count"`
	VSize    int  `json:"vsize"`
	TotalFee int  `json:"totalFee"`
}
