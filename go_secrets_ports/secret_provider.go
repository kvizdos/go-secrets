package go_secrets_ports

import (
	"context"
)

type SecretProvider interface {
	Get(ctx context.Context, key string) (string, error)
}
