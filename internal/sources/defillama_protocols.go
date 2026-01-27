package sources

import (
	"encoding/json"
	"fmt"
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
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://api.llama.fi/protocols")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("defillama returned status %d", resp.StatusCode)
	}

	var protocols []DefiLlamaProtocol
	if err := json.NewDecoder(resp.Body).Decode(&protocols); err != nil {
		return nil, err
	}

	return protocols, nil
}
