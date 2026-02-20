package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
)

type MempoolFees struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
}

type MempoolStats struct {
	Count    int `json:"count"`
	VSize    int `json:"vsize"`
	TotalFee int `json:"total_fee"`
}

func GetBitcoinFees(ctx context.Context) (*MempoolFees, error) {
	var fees MempoolFees
	url := "https://mempool.space/api/v1/fees/recommended"
	if err := sources.FetchJSON(ctx, url, &fees); err != nil {
		return nil, err
	}

	return &fees, nil
}

func GetBitcoinMempool(ctx context.Context) (*MempoolStats, error) {
	var stats MempoolStats
	url := "https://mempool.space/api/mempool"

	if err := sources.FetchJSON(ctx, url, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
