package httpx

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories" // Importação correta do repo
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBitcoinNetworkHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()

	// 1. Ajuste: O handler agora pede o NetworkRepository.
	// Passamos nil porque, no teste de cache, o fluxo deve retornar antes de tocar no repo.
	var mockRepo *repositories.NetworkRepository

	cachedResp := models.BitcoinNetworkResponse{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false, // No cache o valor original é false
		},
		BlockHeight:         100,
		HashrateTHs:         123.4,
		Difficulty:          999,
		AvgBlockTimeSeconds: 600,
		Trend:               models.TrendStable,
	}

	cacheKey := cache.KeyBitcoinNetwork
	cache.Set(c, cacheKey, cachedResp, cache.TTLNetworkStats)

	req := httptest.NewRequest(http.MethodGet, "/bitcoin/network", nil)
	rr := httptest.NewRecorder()

	// 2. ATUALIZADO: Passando o mockRepo
	handler := BitcoinNetworkHandler(c, mockRepo)
	handler(rr, req)

	// Validações
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	// O handler deve marcar "cached":true ao ler do cache
	if !strings.Contains(rr.Body.String(), `"cached":true`) {
		t.Fatalf("expected cached=true, got %s", rr.Body.String())
	}
}
