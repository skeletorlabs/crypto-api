package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
)

type GlobalNodesResponse struct {
	TotalNodes int64 `json:"total_nodes"`
}

func GetGlobalNodesCount(ctx context.Context) (int64, error) {
	url := "https://bitnodes.io/api/v1/snapshots/latest/"

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		// retry once
		resp, err = sources.HttpClient.Do(req)
		if err != nil {
			return 0, sources.ErrUpstreamTimeout
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, sources.ErrUpstreamBadStatus
	}

	var data GlobalNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.TotalNodes, nil
}
