package market

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"net/http"
	"time"
)

func GetProtocols(ctx context.Context) ([]DefillamaProtocol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.llama.fi/protocols", nil)
	if err != nil {
		return nil, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		// retry once
		resp, err = sources.HttpClient.Do(req)
		if err != nil {
			return nil, sources.ErrUpstreamTimeout
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, sources.ErrUpstreamBadStatus
	}

	var protocols []DefillamaProtocol
	if err := json.NewDecoder(resp.Body).Decode(&protocols); err != nil {
		return nil, err
	}

	return protocols, nil
}
