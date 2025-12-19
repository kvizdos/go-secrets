package go_secrets_transformers

import (
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_transformers "github.com/kvizdos/go-secrets/internal/adapters/key_transformers"
)

func NewGenericTransformer(transformFunc func(string) string) go_secrets_ports.Transformer {
	return secret_transformers.NewGenericTransformer(transformFunc)
}

func NewEnvTransformer() go_secrets_ports.Transformer {
	return secret_transformers.NewEnvKeyTransformer()
}
