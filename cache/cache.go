package cache

import (
	"context"
	"errors"
	"github.com/valkey-io/valkey-go"
	"os"
	"strings"
)

type DistributedCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	SetField(ctx context.Context, key string, field string, value string) error
	Delete(ctx context.Context, key string) error
	Close()
}

func NewClient() (valkey.Client, error) {
	envUrls := os.Getenv("VALKEY_URLS")
	if envUrls == "" {
		return nil, errors.New("VALKEY_URLS is not set")
	}
	urls := strings.Split(envUrls, ",")
	return valkey.NewClient(valkey.ClientOption{InitAddress: urls})
}

func New() (DistributedCache, error) {
	cacheClient, err := NewClient()
	return &ParserValkeyCache{
		valkeyCache: cacheClient,
	}, err
}
