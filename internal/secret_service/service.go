package secret_service

import (
	"context"
	"errors"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

var _ go_secrets_ports.SecretService = (*GoSecrets)(nil)

type GoSecrets struct {
	secrets go_secrets_ports.SecretProvider
	config  go_secrets_ports.SecretProvider
}

func (gs *GoSecrets) SetSecretProvider(provider go_secrets_ports.SecretProvider) {
	gs.secrets = provider
}

func (gs *GoSecrets) SetConfigProvider(provider go_secrets_ports.SecretProvider) {
	gs.config = provider
}

func (gs *GoSecrets) Get(ctx context.Context, valueType go_secrets_types.GoSecretType, key string) (string, error) {
	if key == "" {
		return "", go_secrets_types.ErrSecretKeyEmpty
	}

	if valueType == go_secrets_types.CONFIG {
		if gs.config == nil {
			return "", go_secrets_types.ErrProviderNotConfigured
		}
		return gs.config.Get(ctx, key)
	}
	if gs.secrets == nil {
		return "", go_secrets_types.ErrProviderNotConfigured
	}
	return gs.secrets.Get(ctx, key)
}

func (gs *GoSecrets) ExecuteSecret(ctx context.Context, key string, do func(string) error) error {
	if gs.secrets == nil {
		return go_secrets_types.ErrProviderNotConfigured
	}

	// LATER, will replce with Secrets.Do
	value, err := gs.secrets.Get(ctx, key)
	if err != nil {
		return err
	}
	executionErr := do(value)
	if executionErr != nil {
		return errors.Join(err, go_secrets_types.ErrExecutionFail)
	}
	return nil
}
