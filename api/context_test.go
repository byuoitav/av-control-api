package api

import (
	"context"
	"testing"
)

func TestCtxRequestID(t *testing.T) {
	ctx := context.Background()
	id := "12345"

	ctx = WithRequestID(ctx, id)
	got := CtxRequestID(ctx)
	if got != id {
		t.Fatalf("expected %q, got %q", id, got)
	}
}

func TestCtxRequestIDEmpty(t *testing.T) {
	got := CtxRequestID(context.Background())
	if got != "" {
		t.Fatalf("expected an empty id, got %q", got)
	}
}
