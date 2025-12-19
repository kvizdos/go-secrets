package go_secrets

import (
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
	"github.com/kvizdos/go-secrets/internal/secret_service"
)

type GoSecretsOption func(go_secrets_ports.SecretService)

func New(option ...GoSecretsOption) go_secrets_ports.SecretServiceGetter {
	out := new(secret_service.GoSecrets)
	for _, opt := range option {
		opt(out)
	}
	return out
}

func WithSecretProvider(provider go_secrets_ports.SecretProvider) GoSecretsOption {
	return func(gs go_secrets_ports.SecretService) {
		gs.RegisterChannel(go_secrets_types.Channel_Secrets, provider)
	}
}

func WithConfigProvider(provider go_secrets_ports.SecretProvider) GoSecretsOption {
	return func(gs go_secrets_ports.SecretService) {
		gs.RegisterChannel(go_secrets_types.Channel_Config, provider)
	}
}

func WithCustomChannel(channel go_secrets_types.Channel, provider go_secrets_ports.SecretProvider) GoSecretsOption {
	return func(gs go_secrets_ports.SecretService) {
		gs.RegisterChannel(channel, provider)
	}
}
