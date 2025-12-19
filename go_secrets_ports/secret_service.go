package go_secrets_ports

import (
	"context"

	"github.com/kvizdos/go-secrets/go_secrets_types"
)

type SecretService interface {
	SetSecretProvider(provider SecretProvider)
	SetConfigProvider(provider SecretProvider)

	SecretServiceGetter
	SecretServiceExecutor
}

type SecretServiceGetter interface {
	Get(ctx context.Context, valueType go_secrets_types.GoSecretType, key string) (string, error)
}

type SecretServiceExecutor interface {
	ExecuteSecret(ctx context.Context, key string, do func(string) error) error
}
