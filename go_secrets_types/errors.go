package go_secrets_types

import (
	"errors"
)

var (
	ErrSecretKeyEmpty        = errors.New("secret key is empty")
	ErrSecretNotFound        = errors.New("secret not found")
	ErrAccessDenied          = errors.New("access denied")
	ErrLookupFailed          = errors.New("lookup failed")
	ErrExecutionFail         = errors.New("secret execution error")
	ErrProviderNotConfigured = errors.New("provider not configured")

	ErrSecretIDInvalid = errors.New("secret ID is invalid")
)
