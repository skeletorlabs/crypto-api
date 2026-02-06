package market

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CoinGeckoPriceResponse map[string]map[string]float64

func GetPriceUSD(ctx context.Context, token string) (float64, error) {
	token = strings.ToLower(token)

	url := fmt.Sprintf(
		"https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd",
		token,
	)

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

	var data CoinGeckoPriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	price, ok := data[token]["usd"]
	if !ok {
		return 0, fmt.Errorf("price not found for %s", token)
	}

	return price, nil
}
