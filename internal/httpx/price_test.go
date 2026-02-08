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

func TestPriceHandler_MissingToken(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := PriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/price/", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPriceHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()
	const fakePrice = 123.45
	cacheKey := cache.KeyMarketPrice("bitcoin")

	cache.Set[models.PriceResponse](c, cacheKey, models.PriceResponse{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false, // O handler vai mudar para true ao ler do cache
		},
		Token: "bitcoin",
		USD:   fakePrice,
	}, time.Minute)

	handler := PriceHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/price/bitcoin", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"cached":true`) {
		t.Fatalf("expected cached=true, got %s", body)
	}

	if !strings.Contains(body, "123.45") {
		t.Fatalf("expected price 123.45, got %s", body)
	}
}

func TestPriceHandler_UpstreamError(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := PriceHandler(c)

	// invalid token to trigger upstream error
	req := httptest.NewRequest(http.MethodGet, "/price/invalidToken", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code == http.StatusOK {
		t.Fatalf("expected error code, got 200")
	}
}
