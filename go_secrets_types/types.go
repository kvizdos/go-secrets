package go_secrets_types

type GoSecretType string

const (
	SECRET GoSecretType = "secret"
	CONFIG GoSecretType = "config"
)
