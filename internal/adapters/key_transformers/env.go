package secret_transformers

import (
	"fmt"
	"os"
)

type envKeyTransformer struct {
	get func(string) string
}

/*
 * EnvKeyTransformer will take a key and return the value from the environment.
 */
func NewEnvKeyTransformer() *envKeyTransformer {
	return &envKeyTransformer{
		get: os.Getenv,
	}
}

func NewTestRewriter(get func(string) string) *envKeyTransformer {
	return &envKeyTransformer{
		get: get,
	}
}

func (r *envKeyTransformer) Transform(key string) string {
	realKey := r.get(key)

	if realKey == "" {
		panic(fmt.Errorf("Missing key transform for '%s'", key))
	}

	return realKey
}
