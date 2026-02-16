package bitcoin

import (
	"context"
	"crypto-api/internal/sources"
	"time"
)

type GlobalNodesResponse struct {
	TotalNodes int64 `json:"total_nodes"`
}

func GetGlobalNodesCount(ctx context.Context) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var data GlobalNodesResponse
	url := "https://bitnodes.io/api/v1/snapshots/latest/"

	if err := sources.FetchJSON(ctx, url, &data); err != nil {
		return 0, err
	}

	return data.TotalNodes, nil
}
