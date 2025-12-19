package main

import (
	"context"
	"fmt"
	"os"

	go_secrets "github.com/kvizdos/go-secrets"
	"github.com/kvizdos/go-secrets/go_secrets_providers"
	"github.com/kvizdos/go-secrets/go_secrets_transformers"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

func main() {
	chained := go_secrets_transformers.ChainTransformers(
		go_secrets_transformers.NewEnvTransformer(),
		go_secrets_transformers.NewGenericTransformer(func(s string) string {
			return fmt.Sprintf("/%s/%s", os.Getenv("ENV"), s)
		}),
	)

	provider := go_secrets_providers.NewTestingProvider(func(s string) string {
		switch s {
		case "/env/hello":
			return "AJHH"
		default:
			return "???"
		}
	})

	secrets := go_secrets.New(
		go_secrets.WithConfigProvider(provider),
		go_secrets.WithTransformer(go_secrets_types.Channel_Config, chained),
	)

	out, err := secrets.Get(context.Background(), go_secrets_types.Channel_Config, "AHAHAHA")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(out)
}
