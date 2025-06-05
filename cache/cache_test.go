package cache

import (
	"os"
	"testing"
)

func TestNewSuccess(t *testing.T) {
	//Note: this test requires your valkey docker stack to be up
	_ = os.Setenv("VALKEY_URLS", "localhost:6379")
	cache, err := New()
	if err != nil {
		t.Error(err)
	}
	if cache == nil {
		t.Error("Expected cache client")
	}
	cache.Close()
}

func TestNewFail(t *testing.T) {
	_ = os.Setenv("VALKEY_URLS", "")
	_, err := New()
	if err == nil {
		t.Error("Expected error when VALKEY_URLS is empty")
	}
}
