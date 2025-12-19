package secret_providers_aws_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_aws "github.com/kvizdos/go-secrets/internal/adapters/providers/aws"
	secret_providers_test "github.com/kvizdos/go-secrets/internal/adapters/providers/test"
)

func TestSSM_Contract_Integration(t *testing.T) {
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

	ssmClient := ssm.NewFromConfig(cfg, func(o *ssm.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	seedSSM(t, ctx, ssmClient)

	secret_providers_test.AssertSecretProviderContract(t, true, func() go_secrets_ports.SecretProvider {
		return secret_providers_aws.NewSSMWithClient(ssmClient)
	})
}

func seedSSM(t *testing.T, ctx context.Context, c *ssm.Client) {
	t.Helper()

	_, err := c.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String("/key"),
		Type:      types.ParameterTypeSecureString,
		Value:     aws.String("secret"),
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		t.Fatalf("put /key: %v", err)
	}

}
