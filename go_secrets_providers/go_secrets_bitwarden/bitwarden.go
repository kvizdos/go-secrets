package go_secrets_bitwarden

import (
	"github.com/bitwarden/sdk-go"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_managers "github.com/kvizdos/go-secrets/internal/adapters/providers/managers"
)

func New(bwClient sdk.BitwardenClientInterface) go_secrets_ports.SecretProvider {
	return secret_providers_managers.NewBitwardenProvider(bwClient.Secrets().Get)
}
