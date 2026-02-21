package intelligence

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"encoding/json"
	"log"
	"net/http"
)

// IntelligenceSnapshotHandler delivers the latest persisted intelligence snapshot.
func IntelligenceSnapshotHandler(c *cache.MemoryCache, repo *repositories.IntelligenceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cachedKey := cache.KeyIntelligenceLatestSnapshot

		// 1. Cache first
		if snapshot, ok := cache.Get[models.IntelligenceSnapshot](c, cachedKey); ok {
			if err := json.NewEncoder(w).Encode(struct {
				Meta models.Meta                 `json:"meta"`
				Data models.IntelligenceSnapshot `json:"data"`
			}{
				Meta: models.Meta{UpdatedAt: snapshot.CreatedAt, Cached: true},
				Data: snapshot,
			}); err != nil {
				log.Printf("[http] failed to encode cached intelligence snapshot: %v", err)
			}
			return
		}

		// 2. Database fallback
		snapshot, err := repo.GetLatest(r.Context())
		if err != nil {
			http.Error(w, "Intelligence data unavailable", http.StatusServiceUnavailable)
			return
		}

		// 3. Warm cache
		cache.Set(c, cachedKey, *snapshot, cache.TTLIntelligenceSnapshot)

		if err := json.NewEncoder(w).Encode(struct {
			Meta models.Meta                  `json:"meta"`
			Data *models.IntelligenceSnapshot `json:"data"`
		}{
			Meta: models.Meta{UpdatedAt: snapshot.CreatedAt, Cached: false},
			Data: snapshot,
		}); err != nil {
			log.Printf("[http] failed to encode intelligence snapshot: %v", err)
		}
	}
}
