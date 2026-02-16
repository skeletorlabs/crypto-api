package market

import (
	"context"
	"crypto-api/internal/engine/bitcoin/correlation"
	"crypto-api/internal/sources"
	"fmt"
	"time"
)

const (
	llamaBaseUrl = "https://api.llama.fi"
	coinsBaseUrl = "https://coins.llama.fi"
)

func GetChains(ctx context.Context) ([]DefillamaChain, error) {
	url := fmt.Sprintf("%s/chains", llamaBaseUrl)

	var chains []DefillamaChain

	if err := sources.FetchJSON(ctx, url, &chains); err != nil {
		return nil, err
	}

	return chains, nil
}

func GetBTCPriceHistory(ctx context.Context, days int) ([]correlation.DataPoint, error) {
	startTime := time.Now().AddDate(0, 0, -days).Unix()

	url := fmt.Sprintf("%s/chart/coingecko:bitcoin?start=%d&span=%d&period=1d",
		coinsBaseUrl, startTime, days)

	var data DefillamaChart
	if err := sources.FetchJSON(ctx, url, &data); err != nil {
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
