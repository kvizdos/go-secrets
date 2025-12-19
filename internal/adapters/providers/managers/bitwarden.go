package secret_providers_managers

import (
	"context"
	"errors"
	"strings"

	sdk "github.com/bitwarden/sdk-go"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

type bwGetter func(secretID string) (*sdk.SecretResponse, error)

type bitwardenProvider struct {
	read bwGetter
	// Implement BitwardenManager methods here
}

func NewBitwardenProvider(reader bwGetter) *bitwardenProvider {
	return &bitwardenProvider{
		read: reader,
	}
}

func (b *bitwardenProvider) Get(ctx context.Context, key string) (string, error) {
	resp, err := b.read(key)
	if err != nil {
		str := err.Error()

		if strings.Contains(str, "Invalid command value") {
			return "", errors.Join(err, go_secrets_types.ErrSecretIDInvalid, go_secrets_types.ErrSecretNotFound)
		}

		if strings.Contains(str, "Resource not found.") {
			return "", errors.Join(err, go_secrets_types.ErrSecretNotFound)
		}

		return "", err
	}

	return resp.Value, nil
}
