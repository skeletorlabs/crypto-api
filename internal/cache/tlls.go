package cache

import "time"

const (
	// Precomputed intelligence snapshots
	TTLIntelligenceSnapshot = 15 * time.Minute

	// Network statistics (block data, halving state)
	TTLNetworkStats = 15 * time.Minute

	// Macro data (e.g. M2 supply from repository)
	TTLMacroData = 24 * time.Hour

	// Market price cache
	TTLBitcoinPrice = 1 * time.Minute

	// Bitcoin fee estimates
	TTLBitcoinFees = 1 * time.Minute
)
