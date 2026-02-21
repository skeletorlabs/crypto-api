package bitcoin

import (
	"crypto-api/internal/cache"
	"crypto-api/internal/models"
	"crypto-api/internal/storage/repositories"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBitcoinNetworkHandler_FromCache(t *testing.T) {
	c := cache.NewMemoryCache()

	var mockRepo *repositories.NetworkRepository

	cachedResp := models.BitcoinNetworkResponse{
		Meta: models.Meta{
			UpdatedAt: time.Now().UTC(),
			Cached:    false,
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

	handler := BitcoinNetworkHandler(c, mockRepo)
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), `"cached":true`) {
		t.Fatalf("expected cached=true, got %s", rr.Body.String())
	}
}
