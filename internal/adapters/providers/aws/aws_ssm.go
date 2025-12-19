package secret_providers_aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"github.com/kvizdos/go-secrets/go_secrets_types"
)

var _ go_secrets_ports.SecretProvider = (*ssmProvider)(nil)

type SSMGetter interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

// ssmProvider reads secrets from AWS SSM Parameter Store.
// Keys should be full parameter names (e.g., "/prod/app/db_password").
// It reads with decryption enabled (so it works for SecureString too).
type ssmProvider struct {
	ssm SSMGetter
}

func NewSSM(cfg aws.Config) *ssmProvider {
	return &ssmProvider{ssm: ssm.NewFromConfig(cfg)}
}

// NewWithClient is handy for tests or when you already have a configured client.
func NewSSMWithClient(client SSMGetter) *ssmProvider {
	return &ssmProvider{ssm: client}
}

func (p *ssmProvider) Get(ctx context.Context, key string) (string, error) {
	out, err := p.ssm.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		// Make common cases clearer while preserving the original error.
		var nf *types.ParameterNotFound
		if errors.As(err, &nf) {
			return "", go_secrets_types.ErrSecretNotFound
		}
		var ad *types.AccessDeniedException
		if errors.As(err, &ad) {
			return "", go_secrets_types.ErrAccessDenied
		}
		return "", errors.Join(fmt.Errorf("ssm: get parameter %s: %w", key, err), go_secrets_types.ErrLookupFailed)
	}

	if out.Parameter == nil || out.Parameter.Value == nil {
		return "", go_secrets_types.ErrSecretNotFound
	}
	return aws.ToString(out.Parameter.Value), nil
}
