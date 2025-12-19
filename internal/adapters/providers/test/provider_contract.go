package secret_providers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

func AssertSecretProviderContract(
	t *testing.T,
	isIntegration bool,
	p func() go_secrets_ports.SecretProvider,
) {
	t.Helper()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		val, err := p().Get(ctx, "/key")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if val != "secret" {
			t.Fatalf("expected secret, got %q", val)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := p().Get(ctx, "/missing")
		if !errors.Is(err, go_secrets_types.ErrSecretNotFound) {
			t.Fatalf("expected ErrSecretNotFound, got %v", err)
		}
	})

	if !isIntegration {
		t.Run("access denied", func(t *testing.T) {
			_, err := p().Get(ctx, "/denied")
			if !errors.Is(err, go_secrets_types.ErrAccessDenied) {
				t.Fatalf("expected ErrAccessDenied, got %v", err)
			}
		})

		t.Run("unknown error", func(t *testing.T) {
			_, err := p().Get(ctx, "/boom")
			if !errors.Is(err, go_secrets_types.ErrLookupFailed) {
				t.Fatalf("expected ErrLookupFailed, got %v", err)
			}
		})
	}
}
