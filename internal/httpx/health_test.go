package httpx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	HealthHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	body := rec.Body.String()
	if body != `{"status":"ok"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}
