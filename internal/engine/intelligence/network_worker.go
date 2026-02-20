package intelligence

import (
	"context"
	"crypto-api/internal/cache"
	"crypto-api/internal/engine/bitcoin/halving"
	"crypto-api/internal/engine/bitcoin/valuation"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/bitcoin"
	"crypto-api/internal/storage/repositories"
	"log"
	"time"
)

// StartNetworkWorker periodically refreshes Bitcoin network data.
func StartNetworkWorker(ctx context.Context, c *cache.MemoryCache, repo *repositories.NetworkRepository) {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		// Run once at startup
		refreshNetwork(ctx, c, repo)

		for {
			select {
			case <-ticker.C:
				refreshNetwork(ctx, c, repo)

			case <-ctx.Done():
				return
			}
		}
	}()
}

func refreshNetwork(ctx context.Context, c *cache.MemoryCache, repo *repositories.NetworkRepository) {
	data, err := bitcoin.GetBitcoinNetwork(ctx)
	if err != nil {
		log.Printf("[network-worker] fetch error: %v", err)
		return
	}

	// Compute trend based on previous stored value
	prev, err := repo.GetLatest(ctx)
	if err == nil && prev != nil {
		if prev.BlockHeight == data.BlockHeight {
			// No new block, skip persistence
			return
		}
	}

	hasPrev := err == nil && prev != nil
	prevAvg := 0.0
	if hasPrev {
		prevAvg = prev.AvgBlockTimeSeconds
	}

	trendStatus := valuation.CalculateTrend(
		data.AvgBlockTime,
		prevAvg,
		hasPrev,
	)

	now := time.Now().UTC()
	halvingState := halving.Compute(int(data.BlockHeight), data.AvgBlockTime/60, now)

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

	if err := repo.Save(ctx, resp); err != nil {
		log.Printf("[network-worker] save error: %v", err)
	}

	_ = repo.DeleteOlderThan(ctx, 90)

	cache.Set(c, cache.KeyBitcoinNetwork, resp, cache.TTLNetworkStats)
}
