package go_secrets_providers_test

import (
	"context"
	"testing"
	"time"

	"github.com/kvizdos/go-secrets/go_secrets_providers"
)

type fakeProvider struct {
	calls int
	val   string
}

func (f *fakeProvider) Get(ctx context.Context, key string) (string, error) {
	f.calls++
	return f.val, nil
}

func TestWithPreflights_Order(t *testing.T) {
	base := &fakeProvider{val: "secret"}

	p := go_secrets_providers.WithPreflights(
		base,
		go_secrets_providers.WithCacheTTL(time.Minute),
		go_secrets_providers.WithSingleFlight(),
	)

	_, _ = p.Get(context.Background(), "k")
	_, _ = p.Get(context.Background(), "k")

	if base.calls != 1 {
		t.Fatalf("expected cache to prevent second call")
	}
}
