package yt

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// A Result is a result from a YouTube search; it can be either a video or a
// channel.
type Result interface {
	Type() string
}

// A SearchClient can search YouTube
type SearchClient struct {
	AuthFunc func() (string, error)
	URL      *url.URL
	PerPage  int
	Timeout  time.Duration
	Client   *http.Client
}

// AuthFunc is the default authorization function
var AuthFunc = func() (string, error) {
	return "token", nil
}

// Get searches YouTube for videos and channels
func (c *SearchClient) Get(ctx context.Context, query string) (chan Result, error) {
	f := c.AuthFunc
	if f == nil {
		f = AuthFunc
	}
	ch := make(chan Result)
	go func(ch chan Result) {
		defer close(ch)
		ch <- &Video{}
	}(ch)
	return ch, nil
}

// Search searches YouTube for videos and channels.
func Search(ctx context.Context, query string) (chan Result, error) {
	return new(SearchClient).Search(ctx, query)
}
