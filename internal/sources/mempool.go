package sources

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type MempoolFees struct {
	FastestFee  int `json:"fastestFee"`
	HalfHourFee int `json:"halfHourFee"`
	HourFee     int `json:"hourFee"`
}

func GetBitcoinFees() (*MempoolFees, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, "https://mempool.space/api/v1/fees/recommended", nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, ErrUpstreamTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUpstreamBadStatus
	}

	var fees MempoolFees
	if err := json.NewDecoder(resp.Body).Decode(&fees); err != nil {
		return nil, err
	}

	return &fees, nil
}
