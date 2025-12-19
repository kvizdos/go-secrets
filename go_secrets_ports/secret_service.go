package go_secrets_ports

import (
	"context"

	"github.com/kvizdos/go-secrets/go_secrets_types"
)

type SecretService interface {
	RegisterChannel(channelName go_secrets_types.Channel, provider SecretProvider)

	SecretServiceGetter
	SecretServiceExecutor
}

type SecretServiceGetter interface {
	Get(ctx context.Context, channel go_secrets_types.Channel, key string) (string, error)
}

type SecretServiceExecutor interface {
	ExecuteSecret(ctx context.Context, channel go_secrets_types.Channel, key string, do func(string) error) error
}
