package sources

import (
	"net/http"
	"time"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

var httpClient HTTPDoer = &http.Client{
	Timeout: 10 * time.Second,
}
