package yt

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetInfo(t *testing.T) {
	withInfoServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ContentTypeXWWWFormURLEncoded)
		fmt.Fprint(w, `player_response={"videoDetails":{"videoId": "abcdefghij"}}`)
	}), func() {
		i, err := GetInfo("abcdefghijk")
		if err != nil {
			t.Fatalf("expected no error, got %e", err)
		}
		if i.VideoDetails.ID != "abcdefghij" {
			t.Errorf("expected info.VideoDetails.ID to be %q, got %q", "abcdefghij", i.VideoDetails.ID)
		}
	})
}

func withInfoServer(h http.Handler, f func()) {
	ts := httptest.NewServer(h)
	u, _ := url.Parse(ts.URL)
	defer func(u *url.URL) { InfoURL = u }(InfoURL)
	InfoURL = u
	f()
}
