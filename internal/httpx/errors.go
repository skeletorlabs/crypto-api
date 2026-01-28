package httpx

import (
	"errors"
	"net/http"
	"strings"

	"crypto-api/internal/sources"
)

type HTTPError struct {
	Status  int
	Message string
}

func MapError(err error) HTTPError {
	if err == nil {
		return HTTPError{Status: 200, Message: ""}
	}

	switch {
	case errors.Is(err, sources.ErrUpstreamTimeout):
		return HTTPError{Status: http.StatusGatewayTimeout, Message: "Upstream timeout"}

	case errors.Is(err, sources.ErrUpstreamBadStatus):
		return HTTPError{Status: http.StatusBadGateway, Message: "Upstream error"}

	case strings.Contains(err.Error(), "price not found"):
		return HTTPError{Status: http.StatusBadRequest, Message: err.Error()}
	default:
		return HTTPError{Status: http.StatusInternalServerError, Message: "Internal error"}
	}
}
