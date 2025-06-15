package validator

import "context"

type MockCache struct {
}

func (p *MockCache) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (p *MockCache) Set(ctx context.Context, key, value string) error {
	return nil
}

func (p *MockCache) SetField(ctx context.Context, key, field, value string) error {
	return nil
}

func (p *MockCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (p *MockCache) Close() {
}
