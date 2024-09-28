package runtime

import (
	"context"

	"github.com/tetratelabs/wazero/api"
)

type GuestModule struct {
	api.Module
	ReturnDataPtr uint32
}

// NewGuestModule creates a new GuestModule.
func NewGuestModule(ctx context.Context, m api.Module) GuestModule {
	// get a memory buffer for return values, if possible
	var returnDataPtr uint32
	malloc := m.ExportedFunction("malloc")
	if malloc != nil {
		res, err := malloc.Call(ctx, 256)
		if err == nil {
			returnDataPtr = uint32(res[0])
		}
	}

	return GuestModule{
		Module:        m,
		ReturnDataPtr: returnDataPtr,
	}
}
