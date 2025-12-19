package secret_providers

import (
	"context"
	"os"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

var _ go_secrets_ports.SecretProvider = (*envSecretProvider)(nil)

type envSecretProvider struct {
	getEnv func(string) string
}

func NewEnvSecretProvider(getEnv func(string) string) *envSecretProvider {
	if getEnv == nil {
		getEnv = os.Getenv
	}
	return &envSecretProvider{
		getEnv: getEnv,
	}
}

func (p *envSecretProvider) Get(ctx context.Context, key string) (string, error) {
	val := p.getEnv(key)
	if val == "" {
		return "", go_secrets_types.ErrSecretNotFound
	}
	return val, nil
}
