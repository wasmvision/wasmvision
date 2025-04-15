package storage

import (
	"os"
	"testing"

	"github.com/go-redis/redismock/v9"
)

func TestRedisStorage(t *testing.T) {

	client, mock := redismock.NewClientMock()

	mock.ExpectHSet("proc", "key", "value").SetVal(0)
	mock.ExpectHGet("proc", "key").SetVal("value")

	// just a valid redis url. no need for a local
	// redis instance.
	os.Setenv("", "redis://127.0.0.1:6379/1")

	s := NewRedisStorage()
	s.db = client

	err := s.Set("proc", "key", "value")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	val, ok := s.Get("proc", "key")
	if !ok {
		t.Errorf("key not found")
	}

	if string(val) != "value" {
		t.Errorf("unexpected value: %s", val)
	}

}
