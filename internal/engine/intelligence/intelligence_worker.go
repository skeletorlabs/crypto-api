package intelligence

import (
	"context"
	"crypto-api/internal/engine/bitcoin/providers"
	"log"
	"time"
)

func StartDailyWorker(ctx context.Context, provider *providers.IntelligenceProvider) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[daily-worker] panic recovered: %v", r)
			}
		}()

		// Run immediately on startup
		if err := provider.GenerateFullSnapshot(ctx); err != nil {
			log.Printf("[daily-worker] initial snapshot error: %v", err)
		}

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				log.Println("[daily-worker] running scheduled snapshot...")
				if err := provider.GenerateFullSnapshot(ctx); err != nil {
					log.Printf("[daily-worker] scheduled snapshot error: %v", err)
				}
			}
		}
	}()
}
