package macro

import (
	"context"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/sources"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type FredResponse struct {
	Observations []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"observations"`
}

func fetchFredData(ctx context.Context, limit int) (*FredResponse, error) {
	apiKey := os.Getenv("FRED_API_KEY")
	url := fmt.Sprintf(
		"https://api.stlouisfed.org/fred/series/observations?series_id=WM2NS&api_key=%s&file_type=json&sort_order=desc&limit=%d",
		apiKey, limit,
	)

	// We can centralize the timeout here too
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

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

	var data FredResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func GetM2Supply(ctx context.Context) (float64, time.Time, error) {
	data, err := fetchFredData(ctx, 1)
	if err != nil {
		return 0, time.Time{}, err
	}

	if len(data.Observations) == 0 {
		return 0, time.Time{}, fmt.Errorf("no data found")
	}

	val, _ := strconv.ParseFloat(data.Observations[0].Value, 64)
	date, _ := time.Parse("2006-01-02", data.Observations[0].Date)

	return val, date, nil
}

func GetM2History(ctx context.Context, limit int) ([]correlation.DataPoint, error) {
	data, err := fetchFredData(ctx, limit)
	if err != nil {
		return nil, err
	}

	var history []correlation.DataPoint
	for _, obs := range data.Observations {
		if val, err := strconv.ParseFloat(obs.Value, 64); err == nil {
			date, _ := time.Parse("2006-01-02", obs.Date)
			history = append(history, correlation.DataPoint{
				Date:  date,
				Value: val,
			})
		}
	}

	return history, nil
}
