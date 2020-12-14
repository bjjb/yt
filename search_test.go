package yt

import (
	"context"
	"testing"
)

func TestSearch(t *testing.T) {
	rs, err := Search(context.Background(), "foo")
	if err != nil {
		t.Fatalf("expected no error, got %e", err)
	}
	if cap(rs) != 1 {
		t.Fatalf("expected %d results, got %d", 1, len(rs))
	}
	r := <-rs
	x := r.Type()
	if x != "video" {
		t.Errorf("expected results[0].Type to be %q, got %q", "video", x)
	}
}
