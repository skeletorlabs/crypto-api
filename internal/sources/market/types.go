package market

import (
	"context"
)

// Coingecko types
type CoingeckoPriceResponse map[string]map[string]float64

// Binance types
type BinancePriceResponse struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"` // Binance returns price as a string
}

// Coinbase types
type CoinbasePriceResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

// Kraken types
type KrakenPriceResponse struct {
	Error  []string `json:"error"`
	Result map[string]struct {
		C []string `json:"c"` // Last trade closed [price, lot volume]
	} `json:"result"`
}

type providerStrategy struct {
	URL    func(ticker string) string
	Parser func(ctx context.Context, url string) (float64, error)
}

// Defillama types
type DefillamaChain struct {
	Name        string  `json:"name"`
	TVL         float64 `json:"tvl"`
	TokenSymbol string  `json:"tokenSymbol"`
}

type DefillamaProtocol struct {
	Name     string  `json:"name"`
	Slug     string  `json:"slug"`
	TVL      float64 `json:"tvl"`
	Chain    string  `json:"chain"`
	Category string  `json:"category"`
}

type DefillamaChart struct {
	Coins map[string]struct {
		Symbol     string                `json:"symbol"`
		Confidence float64               `json:"confidence"`
		Prices     []DefillamaPricePoint `json:"prices"`
	} `json:"coins"`
}

type DefillamaPricePoint struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}
