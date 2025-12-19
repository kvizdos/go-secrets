package secret_providers_aws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_aws "github.com/kvizdos/go-secrets/internal/adapters/providers/aws"
	secret_providers_test "github.com/kvizdos/go-secrets/internal/adapters/providers/test"
)

type fakeSSM struct {
	responses map[string]struct {
		out *ssm.GetParameterOutput
		err error
	}
}

func (f *fakeSSM) GetParameter(
	ctx context.Context,
	in *ssm.GetParameterInput,
	_ ...func(*ssm.Options),
) (*ssm.GetParameterOutput, error) {
	r := f.responses[aws.ToString(in.Name)]
	return r.out, r.err
}

func TestSSM_Contract(t *testing.T) {
	secret_providers_test.AssertSecretProviderContract(t, false, func() go_secrets_ports.SecretProvider {
		return secret_providers_aws.NewSSMWithClient(&fakeSSM{
			responses: map[string]struct {
				out *ssm.GetParameterOutput
				err error
			}{
				"/key": {
					out: &ssm.GetParameterOutput{
						Parameter: &types.Parameter{
							Value: aws.String("secret"),
						},
					},
				},
				"/missing": {err: &types.ParameterNotFound{}},
				"/denied":  {err: &types.AccessDeniedException{}},
				"/boom":    {err: errors.New("boom")},
				"/empty": {
					out: &ssm.GetParameterOutput{
						Parameter: &types.Parameter{},
					},
				},
			},
		})
	})
}
