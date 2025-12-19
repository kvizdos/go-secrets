package go_secrets_providers

import (
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers "github.com/kvizdos/go-secrets/internal/adapters/providers"
)

func NewEnvProvider() go_secrets_ports.SecretProvider {
	return secret_providers.NewEnvSecretProvider(nil)
}

func NewTestingProvider(getter func(string) string) go_secrets_ports.SecretProvider {
	return secret_providers.NewEnvSecretProvider(getter)
}

func NewTestingMappedProvider(input map[string]string) go_secrets_ports.SecretProvider {
	getter := func(key string) string {
		if v, ok := input[key]; ok {
			return v
		}
		return ""
	}
	return secret_providers.NewEnvSecretProvider(getter)
}
