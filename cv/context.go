package cv

import (
	"log/slog"

	"github.com/wasmvision/wasmvision/config"
	"github.com/wasmvision/wasmvision/datastore"
	"github.com/wasmvision/wasmvision/datastore/storage"
)

// Context is the configuration for the cv package used when each call is made
// from the guest module.
type Context struct {
	ReturnDataPtr  uint32
	ModelsDir      string
	Config         *config.Store
	FrameStore     *datastore.Frames
	ProcessorStore *datastore.Processors
	EnableCUDA     bool
}

func NewContext(modelsDir string, conf *config.Store, datastorage string, enableCUDA bool) *Context {
	var store datastore.DataStorage
	switch datastorage {
	case "memory":
		slog.Info("Using memory datastorage")
		store = storage.NewMemStorage[string]()

	case "boltdb":
		slog.Info("Using BoltDB datastorage")
		store = storage.NewBoltDBStorage()

	case "redis":
		slog.Info("Using Redis datastorage")
		store = storage.NewRedisStorage()

	case "nats":
		store = storage.NewNatsStorage()

	default:
		slog.Info("Using default datastorage (memory)")
		store = storage.NewMemStorage[string]()
	}

	return &Context{
		ModelsDir:      modelsDir,
		Config:         conf,
		FrameStore:     datastore.NewFrames(),
		ProcessorStore: datastore.NewProcessors(store),
		EnableCUDA:     enableCUDA,
	}
}
