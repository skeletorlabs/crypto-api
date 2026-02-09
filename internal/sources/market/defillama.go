package market

import (
	"context"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/sources"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	llamaBaseUrl = "https://api.llama.fi"
	coinsBaseUrl = "https://coins.llama.fi"
)

func GetChains(ctx context.Context) ([]DefillamaChain, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s/chains", llamaBaseUrl)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var chains []DefillamaChain
	if err := json.NewDecoder(resp.Body).Decode(&chains); err != nil {
		return nil, err
	}

	return chains, nil
}

// GetBTCPriceHistory fetches historical BTC prices from DefiLlama
func GetBTCPriceHistory(ctx context.Context, days int) ([]correlation.DataPoint, error) {
	startTime := time.Now().AddDate(0, 0, -days).Unix()

	url := fmt.Sprintf("%s/chart/coingecko:bitcoin?start=%d&span=%d&period=1d",
		coinsBaseUrl, startTime, days)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		return nil, sources.ErrUpstreamTimeout
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, sources.ErrUpstreamBadStatus
	}

	var data DefillamaChart
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	btcData, ok := data.Coins["coingecko:bitcoin"]
	if !ok || len(btcData.Prices) == 0 {
		return nil, fmt.Errorf("no historical data found for bitcoin")
	}

	var history []correlation.DataPoint
	for _, p := range btcData.Prices {
		history = append(history, correlation.DataPoint{
			Date:  time.Unix(p.Timestamp, 0).UTC(),
			Value: p.Price,
		})
	}

	return history, nil
}
