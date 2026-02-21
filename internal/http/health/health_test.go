package health

import (
	"crypto-api/internal/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {

	mockStore := &storage.PostgresStore{}

	handler := HealthHandler(mockStore)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status 200 or 503, got %d", res.StatusCode)
	}

	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] == "" {
		t.Fatal("expected status field in response")
	}

	if _, ok := response["time"]; !ok {
		t.Fatal("expected time field in response")
	}
}
