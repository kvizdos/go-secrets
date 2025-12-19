package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bitwarden/sdk-go"
	go_secrets "github.com/kvizdos/go-secrets"
	"github.com/kvizdos/go-secrets/go_secrets_providers/go_secrets_bitwarden"
	"github.com/kvizdos/go-secrets/go_secrets_transformers"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

func main() {
	// Normally, set outside of Go
	// _ = os.Setenv("A_BITWARDEN_SECRET", "SET_ME")

	env := go_secrets_transformers.NewEnvTransformer()

	apiURL := "https://api.bitwarden.com"
	identityURL := "https://identity.bitwarden.com"

	bitwardenClient, err := sdk.NewBitwardenClient(&apiURL, &identityURL)

	if err != nil {
		panic(err)
	}

	err = bitwardenClient.AccessTokenLogin(os.Getenv("BWS_ACCESS_TOKEN"), nil)
	if err != nil {
		panic(err)
	}

	provider := go_secrets_bitwarden.New(bitwardenClient)

	secrets := go_secrets.New(
		go_secrets.WithConfigProvider(provider),
		go_secrets.WithTransformer(go_secrets_types.Channel_Config, env),
	)

	out, err := secrets.Get(context.Background(), go_secrets_types.Channel_Config, "A_BITWARDEN_SECRET")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(out)
}
