package providers

import (
	"context"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/halving"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/bitcoin"
	"crypto-api/internal/storage"
)

type NetworkProvider struct {
	cache *cache.MemoryCache
	db    *storage.PostgresStore
}

// NewNetworkProvider creates a provider with cache and database access.
func NewNetworkProvider(c *cache.MemoryCache, db *storage.PostgresStore) *NetworkProvider {
	return &NetworkProvider{
		cache: c,
		db:    db,
	}
}

// Update fetches fresh network data, computes derived metrics,
// and updates cache.
func (p *NetworkProvider) Update(ctx context.Context) error {
	// --- FETCH ---
	data, err := bitcoin.GetBitcoinNetwork(ctx)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	// --- TREND ---
	prev, _ := p.db.GetLatestSnapshot(ctx)

	hasPrev := prev != nil
	prevAvg := 0.0
	if hasPrev {
		prevAvg = prev.AvgBlockTime
	}

	trendStatus := valuation.CalculateTrend(
		data.AvgBlockTime,
		prevAvg,
		hasPrev,
	)

	// --- HALVING ---
	halvingState := halving.Compute(
		int(data.BlockHeight),
		data.AvgBlockTime/60,
		now,
	)

	// --- RESPONSE ---
	resp := models.BitcoinNetworkResponse{
		Meta: models.Meta{
			UpdatedAt: now,
			Cached:    false,
		},
		BlockHeight:         data.BlockHeight,
		HashrateTHs:         data.HashrateTHs,
		Difficulty:          data.Difficulty,
		AvgBlockTimeSeconds: data.AvgBlockTime,
		Trend:               trendStatus,
		Halving:             halvingState,
	}

	// --- CACHE ---
	cache.Set(p.cache, cache.KeyBitcoinNetwork, resp, cache.TTLNetworkStats)

	return nil
}
