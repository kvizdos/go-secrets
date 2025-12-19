package secret_transformers

type genericTransformer struct {
	transform func(string) string
}

/*
 * Generic Transformer is insanely powerful.
 */
func NewGenericTransformer(transformer func(string) string) *genericTransformer {
	return &genericTransformer{
		transform: transformer,
	}
}

func (r *genericTransformer) Transform(key string) string {
	return r.transform(key)
}
