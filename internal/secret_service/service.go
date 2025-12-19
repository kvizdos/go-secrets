package secret_service

import (
	"context"
	"errors"
	"sync"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

var _ go_secrets_ports.SecretService = (*GoSecrets)(nil)

type GoSecrets struct {
	mu        sync.Mutex
	providers map[go_secrets_types.Channel]go_secrets_ports.SecretProvider
}

func (gs *GoSecrets) RegisterChannel(channelName go_secrets_types.Channel, provider go_secrets_ports.SecretProvider) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.providers == nil {
		gs.providers = make(map[go_secrets_types.Channel]go_secrets_ports.SecretProvider)
	}
	gs.providers[channelName] = provider
}

func (gs *GoSecrets) Get(ctx context.Context, channel go_secrets_types.Channel, key string) (string, error) {
	if key == "" {
		return "", go_secrets_types.ErrSecretKeyEmpty
	}

	if provider, ok := gs.providers[channel]; ok {
		return provider.Get(ctx, key)
	}

	return "", go_secrets_types.ErrProviderNotConfigured
}

func (gs *GoSecrets) ExecuteSecret(ctx context.Context, channel go_secrets_types.Channel, key string, do func(string) error) error {
	provider, ok := gs.providers[channel]

	if !ok {
		return go_secrets_types.ErrProviderNotConfigured
	}

	// LATER, will replce with Secrets.Do
	value, err := provider.Get(ctx, key)
	if err != nil {
		return err
	}
	executionErr := do(value)
	if executionErr != nil {
		return errors.Join(err, go_secrets_types.ErrExecutionFail)
	}
	return nil
}
