package sources

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type DefiLlamaProtocol struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	TVL      float64 `json:"tvl"`
	Chain    string  `json:"chain"`
	Category string  `json:"category"`
}

func GetProtocols() ([]DefiLlamaProtocol, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.llama.fi/protocols", nil)
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

	var protocols []DefiLlamaProtocol
	if err := json.NewDecoder(resp.Body).Decode(&protocols); err != nil {
		return nil, err
	}

	return protocols, nil
}
