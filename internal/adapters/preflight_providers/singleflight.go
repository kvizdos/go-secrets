package secret_preflights

import (
	"context"

	"github.com/kvizdos/go-secrets/go_secrets_ports"
	"golang.org/x/sync/singleflight"
)

var _ go_secrets_ports.SecretProvider = (*singleflightProvider)(nil)

type singleflightProvider struct {
	base  go_secrets_ports.SecretProvider
	group singleflight.Group
}

func NewSingleflightProvider(base go_secrets_ports.SecretProvider) *singleflightProvider {
	return &singleflightProvider{
		base: base,
	}
}

func (p *singleflightProvider) Get(ctx context.Context, key string) (string, error) {
	ch := p.group.DoChan(key, func() (any, error) {
		return p.base.Get(ctx, key)
	})

	select {
	case res := <-ch:
		if res.Err != nil {
			return "", res.Err
		}
		return res.Val.(string), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
