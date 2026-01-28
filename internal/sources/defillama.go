package sources

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type DefiLlamaChain struct {
	Name        string  `json:"name"`
	TVL         float64 `json:"tvl"`
	TokenSymbol string  `json:"tokenSymbol"`
}

func GetChains() ([]DefiLlamaChain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.llama.fi/chains", nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		// retry once
		resp, err = httpClient.Do(req)
		if err != nil {
			return nil, ErrUpstreamTimeout
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUpstreamBadStatus
	}

	var chains []DefiLlamaChain
	if err := json.NewDecoder(resp.Body).Decode(&chains); err != nil {
		return nil, err
	}

	return chains, nil
}
