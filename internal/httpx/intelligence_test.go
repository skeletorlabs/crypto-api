package httpx

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

func TestIntelligenceHandler(t *testing.T) {
	// 1. Setup de dependências
	// Note: intelRepo aqui deve ser tratado com cuidado se o Handler chamar o DB.
	// Se o Handler usar a interface, um mock seria melhor, mas mantemos conforme sua estrutura.
	intelRepo := &repositories.IntelligenceRepository{}
	intelCache := cache.NewMemoryCache()

	t.Run("Should return data from Cache", func(t *testing.T) {
		// Simula um snapshot já existente no cache
		snapshot := models.IntelligenceSnapshot{
			PriceUSD:    70000.0,
			TrendStatus: "BULLISH",
			CreatedAt:   time.Now().UTC(),
		}

		cache.Set(intelCache, cache.KeyIntelligenceLatestSnapshot, snapshot, cache.TTLIntelligenceSnapshot)

		handler := IntelligenceHandler(intelCache, intelRepo)
		req := httptest.NewRequest(http.MethodGet, "/v1/bitcoin/intelligence", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		// Valida se o JSON retornado tem o campo "cached": true
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
		// Cache novo e vazio
		emptyCache := cache.NewMemoryCache()

		handler := IntelligenceHandler(emptyCache, intelRepo)
		req := httptest.NewRequest(http.MethodGet, "/v1/bitcoin/intelligence", nil)
		rec := httptest.NewRecorder()

		// Como o repo não tem conexão real (é um struct vazio), o Handler deve retornar erro
		handler.ServeHTTP(rec, req)

		// O status 503 ou 404 depende de como seu MapError trata o "not found" do repo
		if rec.Code != http.StatusServiceUnavailable && rec.Code != http.StatusNotFound {
			t.Errorf("expected error status when no data is available, got %d", rec.Code)
		}
	})
}
