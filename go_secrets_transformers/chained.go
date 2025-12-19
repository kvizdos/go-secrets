package go_secrets_transformers

import "github.com/kvizdos/go-secrets/go_secrets_ports"

type transformerChain struct {
	list []go_secrets_ports.Transformer
}

func (c transformerChain) Transform(key string) string {
	for _, t := range c.list {
		key = t.Transform(key)
	}
	return key
}

func ChainTransformers(ts ...go_secrets_ports.Transformer) go_secrets_ports.Transformer {
	return transformerChain{list: ts}
}
