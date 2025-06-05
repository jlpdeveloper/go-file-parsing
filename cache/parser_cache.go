package cache

import (
	"context"
	"github.com/valkey-io/valkey-go"
)

type ParserValkeyCache struct {
	valkeyCache valkey.Client
}

func (p *ParserValkeyCache) Get(ctx context.Context, key string) (string, error) {
	result := p.valkeyCache.Do(ctx, p.valkeyCache.B().Get().Key(key).Build())
	return result.ToString()
}

func (p *ParserValkeyCache) Set(ctx context.Context, key, value string) error {
	return p.valkeyCache.Do(ctx, p.valkeyCache.B().Set().Key(key).Value(value).Build()).Error()
}

func (p *ParserValkeyCache) SetField(ctx context.Context, key, field, value string) error {
	return p.valkeyCache.Do(ctx, p.valkeyCache.B().Hsetnx().Key(key).Field(field).Value(value).Build()).Error()
}

func (p *ParserValkeyCache) Close() {
	p.valkeyCache.Close()
}
