//go:build integration
// +build integration

package secret_providers_aws_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_aws "github.com/kvizdos/go-secrets/internal/adapters/providers/aws"
	secret_providers_test "github.com/kvizdos/go-secrets/internal/adapters/providers/test"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

func TestSecretsManager_Contract_Integration(t *testing.T) {
	ctx := context.Background()

	ls, err := localstack.Run(ctx, "localstack/localstack:latest")
	if err != nil {
		t.Fatalf("start localstack: %v", err)
	}
	t.Cleanup(func() { _ = ls.Terminate(ctx) })

	mappedPort, err := ls.MappedPort(ctx, "4566/tcp")
	if err != nil {
		t.Fatalf("mapped port: %v", err)
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		t.Fatalf("docker provider: %v", err)
	}
	t.Cleanup(func() { _ = provider.Close() })

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		t.Fatalf("daemon host: %v", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		t.Fatalf("aws config: %v", err)
	}

	smClient := secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	seedSecretsManager(t, ctx, smClient)

	// Secrets Manager allows empty strings, but your contract normalizes them away,
	// so pass `true` to enable the empty-value contract cases.
	secret_providers_test.AssertSecretProviderContract(t, true, func() go_secrets_ports.SecretProvider {
		return secret_providers_aws.NewSecretsManagerWithClient(smClient)
	})
}

func seedSecretsManager(t *testing.T, ctx context.Context, c *secretsmanager.Client) {
	t.Helper()

	_, err := c.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String("/key"),
		SecretString: aws.String("secret"),
	})
	if err != nil {
		t.Fatalf("create /key: %v", err)
	}

	// Optional: seed empty secret if you want the empty-value path exercised
	_, err = c.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String("/empty"),
		SecretString: aws.String(""),
	})
	if err != nil {
		t.Fatalf("create /empty: %v", err)
	}
}
