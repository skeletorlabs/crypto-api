package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
)

type HashrateHistoryResponse struct {
	Values []struct {
		X int64   `json:"x"`
		Y float64 `json:"y"`
	} `json:"values"`
}

func GetHashrateHistory(ctx context.Context) (HashrateHistoryResponse, error) {
	url := "https://api.blockchain.info/charts/hash-rate?timespan=30days&format=json&cors=true"

	var data HashrateHistoryResponse
	if err := sources.FetchJSON(ctx, url, &data); err != nil {
		return HashrateHistoryResponse{}, err
	}

	return data, nil
}
