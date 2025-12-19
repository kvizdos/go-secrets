package preflight_tests

import (
	"context"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	secret_preflights "github.com/kvizdos/go-secrets/internal/adapters/preflight_providers"
)

type fakeProvider struct {
	calls int
	val   string
}

func (f *fakeProvider) Get(ctx context.Context, key string) (string, error) {
	f.calls++
	return f.val, nil
}

func TestCache_Hit(t *testing.T) {
	base := &fakeProvider{val: "secret"}
	cache := secret_preflights.NewCacheProvider(base, time.Minute)

	ctx := context.Background()

	v1, _ := cache.Get(ctx, "a")
	v2, _ := cache.Get(ctx, "a")

	if v1 != "secret" || v2 != "secret" {
		t.Fatal("unexpected value")
	}
	if base.calls != 1 {
		t.Fatalf("expected 1 call, got %d", base.calls)
	}
}

func TestCache_Expires(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		base := &fakeProvider{val: "secret"}
		cache := secret_preflights.NewCacheProvider(base, 10*time.Minute)

		ctx := context.Background()

		cache.Get(ctx, "a")
		time.Sleep(5 * time.Minute)
		cache.Get(ctx, "a")

		if base.calls != 1 {
			t.Fatalf("expected 1 pre-expiring, got %d", base.calls)
		}

		time.Sleep(6 * time.Minute)
		cache.Get(ctx, "a")

		if base.calls != 2 {
			t.Fatalf("expected 2 after expiring, got %d", base.calls)
		}
	})

}

func TestCache_Concurrent(t *testing.T) {
	base := &fakeProvider{val: "secret"}
	cache := secret_preflights.NewCacheProvider(base, 10*time.Minute)

	ctx := context.Background()
	var wg sync.WaitGroup

	for range 50 {
		wg.Go(func() {
			cache.Get(ctx, "a")
		})
	}
	wg.Wait()
}
