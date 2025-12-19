package secret_transformers_test

import (
	"fmt"
	"testing"

	secret_transformers "github.com/kvizdos/go-secrets/internal/adapters/key_transformers"
)

func TestGenericTransform(t *testing.T) {

	transformer := secret_transformers.NewGenericTransformer(func(s string) string {
		return fmt.Sprintf("/prod/demo/%s", s)
	})

	transformed := transformer.Transform("test")
	if transformed != "/prod/demo/test" {
		t.Errorf("Expected '/prod/demo/test', got '%s'", transformed)
	}
}
