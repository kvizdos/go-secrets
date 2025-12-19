package secret_transformers_test

import (
	"testing"

	secret_transformers "github.com/kvizdos/go-secrets/internal/adapters/key_transformers"
)

func TestEnvTransform(t *testing.T) {
	demo := func(string) string {
		return "Hello World"
	}
	transformer := secret_transformers.NewTestRewriter(demo)

	transformed := transformer.Transform("test")
	if transformed != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", transformed)
	}
}

func TestEnvTransformMissingKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic, but did not panic")
		}
	}()

	demo := func(string) string {
		return ""
	}
	transformer := secret_transformers.NewTestRewriter(demo)

	transformed := transformer.Transform("test")
	if transformed != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", transformed)
	}
}
