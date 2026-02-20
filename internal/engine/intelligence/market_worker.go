package intelligence

import (
	"context"
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/market"
	"time"
)

// StartMarketWorker periodically fetches market prices
// for all supported assets and updates the intelligence cache.
func StartMarketWorker(ctx context.Context, c *cache.MemoryCache) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				for asset := range SupportedAssets {
					price, err := market.GetPriceUSD(ctx, asset)
					if err != nil {
						continue
					}

					data := models.IntelligencePrice{
						Price:     price,
						CreatedAt: time.Now().UTC(),
					}

					key := cache.KeyIntelligencePrice(asset)
					cache.Set(c, key, data, cache.TTLBitcoinPrice)
				}
			}
		}
	}()
}
