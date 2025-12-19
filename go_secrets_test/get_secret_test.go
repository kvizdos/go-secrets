package go_secrets_test

import (
	"context"
	"errors"
	"testing"

	go_secrets "github.com/kvizdos/go-secrets"
	"github.com/kvizdos/go-secrets/go_secrets_providers"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

func TestGetEmptyKey(t *testing.T) {
	demoEnv := func(key string) string {
		return "Hello World"
	}

	demoProvider := go_secrets_providers.NewTestingProvider(demoEnv)
	secrets := go_secrets.New(
		go_secrets.WithSecretProvider(demoProvider),
		go_secrets.WithConfigProvider(demoProvider),
	)

	for _, channel := range []go_secrets_types.Channel{go_secrets_types.Channel_Secrets, go_secrets_types.Channel_Config} {
		_, err := secrets.Get(context.Background(), channel, "")

		if !errors.Is(err, go_secrets_types.ErrSecretKeyEmpty) {
			t.Errorf("Expected error for empty key")
		}
	}
}

func TestGetFound(t *testing.T) {
	demoEnv := func(key string) string {
		return "Hello World"
	}

	demoProvider := go_secrets_providers.NewTestingProvider(demoEnv)
	secrets := go_secrets.New(
		go_secrets.WithSecretProvider(demoProvider),
		go_secrets.WithConfigProvider(demoProvider),
	)

	for _, channel := range []go_secrets_types.Channel{go_secrets_types.Channel_Secrets, go_secrets_types.Channel_Config} {
		value, err := secrets.Get(context.Background(), channel, "DemoKey")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if value != "Hello World" {
			t.Errorf("Expected value to be 'Hello World', got %v", value)
		}
	}
}

func TestGetNotFound(t *testing.T) {
	demoEnv := func(key string) string {
		return ""
	}

	demoProvider := go_secrets_providers.NewTestingProvider(demoEnv)
	secrets := go_secrets.New(
		go_secrets.WithSecretProvider(demoProvider),
		go_secrets.WithConfigProvider(demoProvider),
	)

	for _, channel := range []go_secrets_types.Channel{go_secrets_types.Channel_Secrets, go_secrets_types.Channel_Config} {
		_, err := secrets.Get(context.Background(), channel, "DemoKey")

		if err == nil {
			t.Fatalf("Expected err!")
		}

		if !errors.Is(err, go_secrets_types.ErrSecretNotFound) {
			t.Fatalf("Got wrong err, expected ErrSecretNotFound, got %v", err)
		}
	}
}

func TestGetNotConfigured(t *testing.T) {
	secrets := go_secrets.New()

	for _, channel := range []go_secrets_types.Channel{go_secrets_types.Channel_Secrets, go_secrets_types.Channel_Config} {
		_, err := secrets.Get(context.Background(), channel, "DemoKey")

		if err == nil {
			t.Fatalf("Expected err!")
		}

		if !errors.Is(err, go_secrets_types.ErrProviderNotConfigured) {
			t.Fatalf("Got wrong err, expected ErrProviderNotConfigured, got %v", err)
		}
	}
}
