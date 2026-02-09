package market

import (
	"bytes"
	"context"
	"crypto-api/internal/sources"
	"errors"
	"io"
	"net/http"
	"testing"
)

type fakeHTTPClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return f.do(req)
}

func TestGetPriceUSD_Success(t *testing.T) {
	originalClient := sources.HttpClient
	defer func() { sources.HttpClient = originalClient }()

	body := `{"bitcoin":{"usd":50000.0}}`

	sources.HttpClient = &fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}, nil
		},
	}

	price, err := GetPriceUSD(context.Background(), "bitcoin")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if price != 50000.0 {
		t.Fatalf("expected price 50000.0, got %f", price)
	}
}

func TestGetPriceUSD_UpstreamError(t *testing.T) {
	originalClient := sources.HttpClient
	defer func() { sources.HttpClient = originalClient }()

	sources.HttpClient = &fakeHTTPClient{
		do: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network down")
		},
	}

	_, err := GetPriceUSD(context.Background(), "bitcoin")
	if !errors.Is(err, sources.ErrUpstreamTimeout) {
		t.Fatalf("expected ErrUpstreamTimeout, got %v", err)
	}
}
