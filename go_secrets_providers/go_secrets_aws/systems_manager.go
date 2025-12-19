package go_secrets_aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_aws "github.com/kvizdos/go-secrets/internal/adapters/providers/aws"
)

func NewSystemsManager(cfg aws.Config) go_secrets_ports.SecretProvider {
	return secret_providers_aws.NewSSM(cfg)
}

// NewWithClient is handy for tests or when you already have a configured client.
func NewSystemsManagerWithClient(client *ssm.Client) go_secrets_ports.SecretProvider {
	return secret_providers_aws.NewSSMWithClient(client)
}
