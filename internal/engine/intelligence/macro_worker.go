package intelligence

import (
	"context"
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/sources/macro"
	"crypto-api/internal/storage/repositories"
	"log"
	"time"
)

// StartMacroWorker periodically synchronizes macroeconomic data.
func StartMacroWorker(ctx context.Context, c *cache.MemoryCache, repo *repositories.MacroRepository) {
	ticker := time.NewTicker(12 * time.Hour)

	go func() {
		// Run once at startup
		executeMacroUpdate(ctx, c, repo)

		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				executeMacroUpdate(ctx, c, repo)
			}
		}
	}()
}

func executeMacroUpdate(ctx context.Context, c *cache.MemoryCache, repo *repositories.MacroRepository) {
	log.Println("[macro-worker] updating M2 supply...")

	m2Value, date, err := macro.GetM2Supply(ctx)
	if err != nil {
		log.Printf("[macro-worker] update failed: %v", err)
		return
	}

	// Persist macro data using source-provided date
	if err := repo.SaveM2(ctx, m2Value, date); err != nil {
		log.Printf("[macro-worker] failed to persist macro data: %v", err)
	}

	resp := models.MacroResponse{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false,
		},
		M2Supply: models.M2Details{
			Value:    m2Value,
			Unit:     "Billions of Dollars",
			DateTime: date,
		},
	}

	cache.Set(c, cache.KeyMacroM2Supply, resp, cache.TTLMacroData)

	log.Println("[macro-worker] macro data updated")
}
