// Package yt provides some functions for working with YouTube in Go libraries
// and applications.
package yt

import (
	"net/http"
)

// DefaultHTTPClient is a default HTTP client
var DefaultHTTPClient *http.Client

func init() {
	DefaultHTTPClient = new(http.Client)
}
