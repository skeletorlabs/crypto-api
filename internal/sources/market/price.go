package market

import (
	"context"
	"crypto-api/internal/sources"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var idToTicker = map[string]string{
	"bitcoin":  "BTC",
	"ethereum": "ETH",
	"solana":   "SOL",
}

var coreApiUrls = map[string]string{
	"binance":  "https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT",
	"coinbase": "https://api.coinbase.com/v2/prices/%s-USD/spot",
	"kraken":   "https://api.kraken.com/0/public/Ticker?pair=%sUSD",
}

var coreProviders = map[string]providerStrategy{
	"binance": {
		URL: func(t string) string {
			return fmt.Sprintf(coreApiUrls["binance"], t)
		},
		Parser: func(ctx context.Context, url string) (float64, error) {
			var res BinancePriceResponse

			if err := sources.FetchJSON(ctx, url, &res); err != nil {
				return 0, err
			}
			return strconv.ParseFloat(res.Price, 64)
		},
	},
	"coinbase": {
		URL: func(t string) string {
			return fmt.Sprintf(coreApiUrls["coinbase"], t)
		},
		Parser: func(ctx context.Context, url string) (float64, error) {
			var res CoinbasePriceResponse
			if err := sources.FetchJSON(ctx, url, &res); err != nil {
				return 0, err
			}
			return strconv.ParseFloat(res.Data.Amount, 64)
		},
	},
	"kraken": {
		URL: func(t string) string {
			if t == "BTC" {
				t = "XBT"
			}
			return fmt.Sprintf(coreApiUrls["kraken"], t)
		},
		Parser: func(ctx context.Context, url string) (float64, error) {
			var res KrakenPriceResponse
			if err := sources.FetchJSON(ctx, url, &res); err != nil {
				return 0, err
			}
			for _, data := range res.Result {
				return strconv.ParseFloat(data.C[0], 64)
			}
			return 0, fmt.Errorf("no price data")
		},
	},
}

func GetPriceUSD(ctx context.Context, token string) (float64, error) {
	token = strings.ToLower(token)
	var ticker string

	if t, ok := idToTicker[token]; ok {
		ticker = t
	} else if len(token) >= 3 && len(token) <= 5 && !strings.Contains(token, " ") {
		ticker = strings.ToUpper(token)
	}

	if ticker != "" {
		priority := []string{"binance", "coinbase", "kraken"}
		for _, provider := range priority {
			strategy, ok := coreProviders[provider]
			if !ok {
				continue
			}

			attemptCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			price, err := strategy.Parser(attemptCtx, strategy.URL(ticker))
			cancel()
			if err == nil {
				return price, nil
			}
		}
	}

	return priceFromCoingecko(ctx, token)
}
