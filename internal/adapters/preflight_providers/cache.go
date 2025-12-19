package secret_preflights

import (
	"context"
	"sync"
	"time"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
)

var _ go_secrets_ports.SecretProvider = (*cacheProvider)(nil)

type cacheEntry struct {
	val       string
	expiresAt time.Time
}

type cacheProvider struct {
	base go_secrets_ports.SecretProvider
	ttl  time.Duration

	mu    sync.RWMutex
	cache map[string]cacheEntry
}

func NewCacheProvider(base go_secrets_ports.SecretProvider, ttl time.Duration) *cacheProvider {
	return &cacheProvider{
		base:  base,
		ttl:   ttl,
		cache: make(map[string]cacheEntry),
	}
}

func (p *cacheProvider) Get(ctx context.Context, key string) (string, error) {
	// Fast path: cached + not expired
	now := time.Now()

	p.mu.RLock()
	if ent, ok := p.cache[key]; ok && now.Before(ent.expiresAt) {
		p.mu.RUnlock()
		return ent.val, nil
	}
	p.mu.RUnlock()

	// Miss/expired: load from base
	val, err := p.base.Get(ctx, key)
	if err != nil {
		return "", err
	}

	// Store
	exp := now.Add(p.ttl)
	p.mu.Lock()
	p.cache[key] = cacheEntry{val: val, expiresAt: exp}
	p.mu.Unlock()

	return val, nil
}
