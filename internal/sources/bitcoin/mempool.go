package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, "https://mempool.space/api/v1/fees/recommended", nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		return nil, sources.ErrUpstreamTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, sources.ErrUpstreamBadStatus
	}

	var fees MempoolFees
	if err := json.NewDecoder(resp.Body).Decode(&fees); err != nil {
		return nil, err
	}

	return &fees, nil
}

func GetBitcoinMempool(ctx context.Context) (*MempoolStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://mempool.space/api/mempool",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		return nil, sources.ErrUpstreamTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, sources.ErrUpstreamBadStatus
	}

	var stats MempoolStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
