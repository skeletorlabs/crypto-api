package macro

import (
	"context"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/sources"
	"fmt"
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

	var data FredResponse
	if err := sources.FetchJSON(ctx, url, &data); err != nil {
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

	obs := data.Observations[0]

	val, err := strconv.ParseFloat(data.Observations[0].Value, 64)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("fred: invalid value format '%s': %w", obs.Value, err)
	}

	date, err := time.Parse("2006-01-02", data.Observations[0].Date)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("fred: invalid date format '%s': %w", obs.Date, err)
	}

	return val, date, nil
}

func GetM2History(ctx context.Context, limit int) ([]correlation.DataPoint, error) {
	data, err := fetchFredData(ctx, limit)
	if err != nil {
		return nil, err
	}

	var history []correlation.DataPoint
	for _, obs := range data.Observations {
		// No histórico, se um ponto falhar (ex: "." em feriados), apenas pulamos
		val, err := strconv.ParseFloat(obs.Value, 64)
		if err != nil {
			continue
		}

		date, err := time.Parse("2006-01-02", obs.Date)
		if err != nil {
			continue
		}

		history = append(history, correlation.DataPoint{
			Date:  date,
			Value: val,
		})
	}

	return history, nil
}
