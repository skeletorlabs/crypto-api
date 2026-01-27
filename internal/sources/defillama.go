package sources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DefiLlamaChain struct {
	Name        string  `json:"name"`
	TVL         float64 `json:"tvl"`
	TokenSymbol string  `json:"tokenSymbol"`
}

func GetChains() ([]DefiLlamaChain, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get("https://api.llama.fi/chains")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("defillama returned status %d", resp.StatusCode)
	}

	var chains []DefiLlamaChain
	if err := json.NewDecoder(resp.Body).Decode(&chains); err != nil {
		return nil, err
	}

	return chains, nil
}
