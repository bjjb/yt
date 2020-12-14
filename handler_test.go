package yt

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	h := new(Handler)
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/abcdefgh123", nil))
	r := w.Result()
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		t.Errorf("expected status code to be %d, got %d", http.StatusOK, r.StatusCode)
	}
}
