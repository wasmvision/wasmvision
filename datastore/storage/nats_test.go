package storage

import (
	"os"
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	natsserver "github.com/nats-io/nats-server/v2/test"
)

func RunDefaultServer() *server.Server {
	opts := natsserver.DefaultTestOptions
	opts.Port = 7030
	opts.Cluster.Name = "testing"
	return natsserver.RunServer(&opts)
}

func TestNATSSet(t *testing.T) {
	s := RunDefaultServer()

	os.Setenv("WASMVISION_STORAGE_NATS_URL", s.ClientURL())
	defer os.Unsetenv("WASMVISION_STORAGE_NATS_URL")

	ds := NewNatsStorage()
	if ds.Err() != nil {
		t.Fatal(ds.Err())
	}
	defer ds.Close()

	err := ds.Set("proc", "key", "value")
	if err != nil {
		t.Error(err)
	}
}
