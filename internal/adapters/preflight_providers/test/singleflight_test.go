package preflight_tests

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	secret_preflights "github.com/kvizdos/go-secrets/internal/adapters/preflight_providers"
)

type blockingProvider struct {
	calls   int32
	unblock chan struct{}
}

func (b *blockingProvider) Get(ctx context.Context, key string) (string, error) {
	atomic.AddInt32(&b.calls, 1)
	<-b.unblock
	return "secret", nil
}

func TestSingleflight_Collapses(t *testing.T) {
	base := &blockingProvider{
		unblock: make(chan struct{}),
	}

	sf := secret_preflights.NewSingleflightProvider(base)
	ctx := context.Background()

	const n = 20
	var wg sync.WaitGroup

	for range n {
		wg.Go(func() {
			val, err := sf.Get(ctx, "k")
			if err != nil || val != "secret" {
				t.Errorf("bad result: %v %v", val, err)
			}
		})
	}

	// ensure all goroutines are waiting inside base.Get
	time.Sleep(10 * time.Millisecond)
	close(base.unblock)
	wg.Wait()

	if atomic.LoadInt32(&base.calls) != 1 {
		t.Fatalf("expected 1 base call, got %d", base.calls)
	}
}

func TestSingleflight_DifferentKeys(t *testing.T) {
	base := &blockingProvider{
		unblock: make(chan struct{}),
	}

	sf := secret_preflights.NewSingleflightProvider(base)
	ctx := context.Background()

	go sf.Get(ctx, "a")
	go sf.Get(ctx, "b")

	time.Sleep(10 * time.Millisecond)
	close(base.unblock)

	time.Sleep(10 * time.Millisecond)

	if atomic.LoadInt32(&base.calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", base.calls)
	}
}

func TestSingleflight_ContextCancel(t *testing.T) {
	base := &blockingProvider{
		unblock: make(chan struct{}),
	}

	sf := secret_preflights.NewSingleflightProvider(base)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := sf.Get(ctx, "k")
	if err == nil {
		t.Fatal("expected context error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled error, got %v", err)
	}
}

func TestSingleflight_ContextDeadlineExceeded(t *testing.T) {
	base := &blockingProvider{
		unblock: make(chan struct{}),
	}

	sf := secret_preflights.NewSingleflightProvider(base)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	_, err := sf.Get(ctx, "k")
	if err == nil {
		t.Fatal("expected context error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded error, got %v", err)
	}
}
