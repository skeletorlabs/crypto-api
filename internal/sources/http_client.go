package sources

import (
	"net/http"
	"time"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

var HttpClient HTTPDoer = &http.Client{
	Timeout: 45 * time.Second,
}
