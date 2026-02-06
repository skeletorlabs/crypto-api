package macro

import (
	"context"
	"crypto-api/internal/sources"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FredResponse struct {
	Observations []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"observations"`
}

func GetM2Supply(ctx context.Context) (string, error) {
	// API key should be handled via env vars in production
	apiKey := "YOUR_API_KEY"
	url := fmt.Sprintf(
		"https://api.stlouisfed.org/fred/series/observations?series_id=WM2NS&api_key=%s&file_type=json&sort_order=desc&limit=1",
		apiKey,
	)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sources.HttpClient.Do(req)
	if err != nil {
		// retry once
		resp, err = sources.HttpClient.Do(req)
		if err != nil {
			return "", sources.ErrUpstreamTimeout
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", sources.ErrUpstreamBadStatus
	}

	var data FredResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Observations) == 0 {
		return "", fmt.Errorf("no data found in FRED response")
	}

	return data.Observations[0].Value, nil
}
