package cache

import "fmt"

// --- Market Keys ---
const (
	KeyMarketChains    = "market:chains:list"
	KeyMarketProtocols = "market:protocols:list"
)

// KeyMarketPrice returns the key for a specific token price
func KeyMarketPrice(token string) string {
	return fmt.Sprintf("market:price:%s", token)
}

// KeyMarketHistory returns the key for a specific token history
func KeyMarketHistory(token string) string {
	return fmt.Sprintf("market:history:%s", token)
}

// --- Bitcoin Keys ---
const (
	KeyBitcoinFees    = "bitcoin:fees:status"
	KeyBitcoinNetwork = "bitcoin:network:status"
	KeyBitcoinMempool = "bitcoin:mempool:status"
)

// --- Macro Keys ---
const (
	KeyMacroM2Supply  = "macro:m2:supply"
	KeyMacroM2History = "macro:m2:history"
)

// --- Intelligence Keys ---
const (
	KeyIntelligenceValuation   = "intelligence:valuation:bitcoin"
	KeyIntelligenceCorrelation = "intelligence:correlation:bitcoin_m2"
)
