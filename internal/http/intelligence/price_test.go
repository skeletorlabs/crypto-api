package intelligence

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//

func TestIntelligencePriceHandler_NotSupported(t *testing.T) {
	c := cache.NewMemoryCache()

	handler := IntelligencePriceHandler(c, nil)

	req := httptest.NewRequest(http.MethodGet, "/intelligence/price/shitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unsupported asset, got %d", rr.Code)
	}
}

func TestIntelligencePriceHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	token := "bitcoin"
	cacheKey := cache.KeyIntelligencePrice(token)

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

	if !strings.Contains(body, "2026-02-18T10:00:00Z") {
		t.Errorf("expected original snapshot timestamp, got %s", body)
	}
	if !strings.Contains(body, "65000") {
		t.Errorf("expected price 65000, got %s", body)
	}
}
