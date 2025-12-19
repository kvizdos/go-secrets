package secret_providers_aws_test

import (
	"context"
	"errors"
	"testing"

	smtypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
	"github.com/kvizdos/go-secrets/go_secrets_ports"
	secret_providers_aws "github.com/kvizdos/go-secrets/internal/adapters/providers/aws"
	secret_providers_test "github.com/kvizdos/go-secrets/internal/adapters/providers/test"
)

type fakeSecretsManager struct {
	responses map[string]struct {
		out *secretsmanager.GetSecretValueOutput
		err error
	}
}

func (f *fakeSecretsManager) GetSecretValue(
	ctx context.Context,
	in *secretsmanager.GetSecretValueInput,
	_ ...func(*secretsmanager.Options),
) (*secretsmanager.GetSecretValueOutput, error) {
	r := f.responses[aws.ToString(in.SecretId)]
	return r.out, r.err
}

type fakeAPIError struct {
	code string
	msg  string
}

func (e *fakeAPIError) Error() string {
	return e.msg
}

func (e *fakeAPIError) ErrorCode() string {
	return e.code
}

func (e *fakeAPIError) ErrorMessage() string {
	return e.msg
}

func (e *fakeAPIError) ErrorFault() smithy.ErrorFault {
	return smithy.FaultClient
}

func TestSecretsManager_Contract(t *testing.T) {
	secret_providers_test.AssertSecretProviderContract(t, false, func() go_secrets_ports.SecretProvider {
		return secret_providers_aws.NewSecretsManagerWithClient(&fakeSecretsManager{
			responses: map[string]struct {
				out *secretsmanager.GetSecretValueOutput
				err error
			}{
				"/key": {
					out: &secretsmanager.GetSecretValueOutput{
						SecretString: aws.String("secret"),
					},
				},
				"/missing": {
					err: &smtypes.ResourceNotFoundException{},
				},
				"/denied": {
					err: &fakeAPIError{
						code: "AccessDeniedException",
						msg:  "access denied",
					},
				},
				"/boom": {
					err: errors.New("boom"),
				},
				"/empty": {
					out: &secretsmanager.GetSecretValueOutput{
						SecretString: aws.String(""),
					},
				},
			},
		})
	})
}
