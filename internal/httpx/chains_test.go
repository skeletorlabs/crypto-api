package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"crypto-api/internal/cache"
)

func TestChainsHandler_OK(t *testing.T) {
	c := cache.NewMemoryCache()
	handler := ChainsHandler(c)

	req := httptest.NewRequest(http.MethodGet, "/chains", nil)
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
