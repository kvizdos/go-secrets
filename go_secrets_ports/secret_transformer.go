package go_secrets_ports

type Transformer interface {
	Transform(string) string
}
