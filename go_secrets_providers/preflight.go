package go_secrets_providers

import (
	"time"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_preflights "github.com/kvizdos/go-secrets/internal/adapters/preflight_providers"
)

func WithPreflights(base go_secrets_ports.SecretProvider, opts ...preflightOpt) go_secrets_ports.SecretProvider {
	for _, opt := range opts {
		base = opt(base)
	}

	return base
}

type preflightOpt func(go_secrets_ports.SecretProvider) go_secrets_ports.SecretProvider

func WithCacheTTL(ttl time.Duration) preflightOpt {
	return func(base go_secrets_ports.SecretProvider) go_secrets_ports.SecretProvider {
		return secret_preflights.NewCacheProvider(base, ttl)
	}
}

func WithSingleFlight() preflightOpt {
	return func(base go_secrets_ports.SecretProvider) go_secrets_ports.SecretProvider {
		return secret_preflights.NewSingleflightProvider(base)
	}
}
