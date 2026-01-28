package sources

import "errors"

var (
	ErrUpstreamTimeout    = errors.New("upstream timeout")
	ErrUpstreamBadStatus = errors.New("upstream bad status")
)
