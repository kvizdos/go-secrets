package secret_providers_aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/aws/smithy-go"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

var _ go_secrets_ports.SecretProvider = (*secretsManagerProvider)(nil)

type SecretsManagerGetter interface {
	GetSecretValue(
		ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.GetSecretValueOutput, error)
}

type secretsManagerProvider struct {
	sm SecretsManagerGetter
}

func NewSecretsManager(cfg aws.Config) *secretsManagerProvider {
	return &secretsManagerProvider{
		sm: secretsmanager.NewFromConfig(cfg),
	}
}

func NewSecretsManagerWithClient(client SecretsManagerGetter) *secretsManagerProvider {
	return &secretsManagerProvider{sm: client}
}

func (p *secretsManagerProvider) Get(ctx context.Context, key string) (string, error) {
	out, err := p.sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})
	if err != nil {
		// Concrete SM error
		var rn *types.ResourceNotFoundException
		if errors.As(err, &rn) {
			return "", go_secrets_types.ErrSecretNotFound
		}

		// Generic AWS API errors (AccessDenied, etc.)
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "AccessDeniedException" {
				return "", go_secrets_types.ErrAccessDenied
			}
		}

		return "", errors.Join(
			fmt.Errorf("secretsmanager: get secret %s: %w", key, err),
			go_secrets_types.ErrLookupFailed,
		)
	}

	// SM returns either SecretString or SecretBinary
	if out.SecretString != nil {
		if *out.SecretString == "" {
			return "", go_secrets_types.ErrSecretNotFound
		}
		return *out.SecretString, nil
	}

	if len(out.SecretBinary) > 0 {
		return "", go_secrets_types.ErrLookupFailed
	}

	return "", go_secrets_types.ErrSecretNotFound
}
