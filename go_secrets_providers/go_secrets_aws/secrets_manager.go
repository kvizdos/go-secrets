package go_secrets_aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_aws "github.com/kvizdos/go-secrets/internal/adapters/providers/aws"
)

func NewSecretsManager(cfg aws.Config) go_secrets_ports.SecretProvider {
	return secret_providers_aws.NewSecretsManager(cfg)
}

func NewSecretsManagerWithClient(client *secretsmanager.Client) go_secrets_ports.SecretProvider {
	return secret_providers_aws.NewSecretsManagerWithClient(client)
}
