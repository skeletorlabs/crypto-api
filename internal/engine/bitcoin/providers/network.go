package providers

import (
	"context"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/halving"
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

	// --- TREND ---
	trendStatus := models.TrendStable

	prev, err := p.db.GetLatestSnapshot(ctx)
	if err == nil && prev != nil {
		diff := data.AvgBlockTime - prev.AvgBlockTime
		const epsilon = 30.0

		if diff > epsilon {
			trendStatus = models.TrendWorsening
		} else if diff < -epsilon {
			trendStatus = models.TrendImproving
		}
	}

	// --- HALVING ---
	halvingState := halving.Compute(
		int(data.BlockHeight),
		data.AvgBlockTime/60, // convert seconds -> minutes
	)

	// --- RESPONSE ---
	resp := models.BitcoinNetworkResponse{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
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
