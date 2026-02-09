package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
)

type HashrateHistoryResponse struct {
	Values []struct {
		X int64   `json:"x"`
		Y float64 `json:"y"`
	} `json:"values"`
}

func GetHashrateHistory(ctx context.Context) (HashrateHistoryResponse, error) {
	url := "https://api.blockchain.info/charts/hash-rate?timespan=30days&format=json&cors=true"

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return HashrateHistoryResponse{}, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		// retry once
		resp, err = sources.HttpClient.Do(req)
		if err != nil {
			return HashrateHistoryResponse{}, sources.ErrUpstreamTimeout
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HashrateHistoryResponse{}, sources.ErrUpstreamBadStatus
	}

	var data HashrateHistoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return HashrateHistoryResponse{}, err
	}

	return data, nil
}
