package cache

import (
	"os"
	"testing"
)

func TestNewSuccess(t *testing.T) {
	//Note: this test requires your valkey docker stack to be up
	_ = os.Setenv("VALKEY_URLS", "localhost:6379")
	_, err := New()
	if err != nil {
		t.Error(err)
	}
}

func TestNewFail(t *testing.T) {
	_ = os.Setenv("VALKEY_URLS", "")
	_, err := New()
	if err == nil {
		t.Error("Error expected for VALKEY_URLS")
	}
}
