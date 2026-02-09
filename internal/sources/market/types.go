package market

// Coingecko types
type CoingeckoPriceResponse map[string]map[string]float64

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
