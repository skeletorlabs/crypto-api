package httpx

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"crypto-api/internal/cache"
	"crypto-api/internal/models"
)

// --- TESTES PARA MARKET PRICE (Hot Fetch) ---

func TestMarketPriceHandler_MissingToken(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := MarketPriceHandler(c)

	// Ajustado para o novo path /market/price/
	req := httptest.NewRequest(http.MethodGet, "/market/price/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestMarketPriceHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	const fakePrice = 123.45
	token := "bitcoin"
	cacheKey := cache.KeyMarketPrice(token)

	// Market cache ainda usa float64 simples
	cache.Set(c, cacheKey, fakePrice, cache.TTLBitcoinPrice)

	handler := MarketPriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/market/price/bitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"cached":true`) {
		t.Errorf("expected cached=true, got %s", body)
	}
	if !strings.Contains(body, "123.45") {
		t.Errorf("expected price 123.45, got %s", body)
	}
}

// --- TESTES PARA INTELLIGENCE PRICE (Blindado) ---

func TestIntelligencePriceHandler_NotSupported(t *testing.T) {
	c := cache.NewMemoryCache()
	// Passamos nil no repo pois o handler deve barrar no IsSupportedAsset antes de consultar o banco
	handler := IntelligencePriceHandler(c, nil)

	req := httptest.NewRequest(http.MethodGet, "/intelligence/price/shitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	// Deve retornar 404 conforme sua implementação de IsSupportedAsset
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unsupported asset, got %d", rr.Code)
	}
}

func TestIntelligencePriceHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	token := "bitcoin"
	cacheKey := cache.KeyIntelligencePrice(token)

	// Criado conforme a sugestão do arquiteto: Struct no cache
	fixedTime := time.Date(2026, 2, 18, 10, 0, 0, 0, time.UTC)
	cache.Set(c, cacheKey, models.IntelligencePrice{
		Price:     65000.0,
		CreatedAt: fixedTime,
	}, cache.TTLBitcoinPrice)

	handler := IntelligencePriceHandler(c, nil)

	req := httptest.NewRequest(http.MethodGet, "/intelligence/price/bitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	// Verifica se o updatedAt retornado é o CreatedAt que estava no cache
	if !strings.Contains(body, "2026-02-18T10:00:00Z") {
		t.Errorf("expected original snapshot timestamp, got %s", body)
	}
	if !strings.Contains(body, "65000") {
		t.Errorf("expected price 65000, got %s", body)
	}
}
