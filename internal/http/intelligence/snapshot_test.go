package intelligence

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSnapshotHandler(t *testing.T) {

	intelRepo := &repositories.IntelligenceRepository{}
	intelCache := cache.NewMemoryCache()

	t.Run("Should return data from Cache", func(t *testing.T) {

		snapshot := models.IntelligenceSnapshot{
			PriceUSD:    70000.0,
			TrendStatus: "BULLISH",
			CreatedAt:   time.Now().UTC(),
		}

		cache.Set(intelCache, cache.KeyIntelligenceLatestSnapshot, snapshot, cache.TTLIntelligenceSnapshot)

		handler := IntelligenceSnapshotHandler(intelCache, intelRepo)
		req := httptest.NewRequest(http.MethodGet, "/v1/bitcoin/intelligence", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var res models.StandardResponse[models.IntelligenceSnapshot]
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode: %v", err)
		}

		if !res.Meta.Cached {
			t.Error("expected Meta.Cached to be true when data is in cache")
		}

		if res.Data.PriceUSD != 70000.0 {
			t.Errorf("expected price 70000, got %f", res.Data.PriceUSD)
		}
	})

	t.Run("Should fail when Cache and DB are empty", func(t *testing.T) {

		emptyCache := cache.NewMemoryCache()

		handler := IntelligenceSnapshotHandler(emptyCache, intelRepo)
		req := httptest.NewRequest(http.MethodGet, "/v1/bitcoin/intelligence", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable && rec.Code != http.StatusNotFound {
			t.Errorf("expected error status when no data is available, got %d", rec.Code)
		}
	})
}
