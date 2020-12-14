package yt

import (
	"net/http"
)

// A Handler is a http.Handler which accepts GET requests for application/json
// on its root, where the path matches a video ID, fetches the response from
// its upstream URL, parses it, and returns it as JSON.
type Handler struct {
	InfoClient      *http.Client
	SearchClient    *http.Client
	StreamingClient *http.Client
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
