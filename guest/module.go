package guest

import (
	"context"

	"github.com/tetratelabs/wazero/api"
)

type Module struct {
	api.Module
	ReturnDataPtr uint32
}

// NewModule creates a new GuestModule.
func NewModule(ctx context.Context, m api.Module) Module {
	// get a memory buffer for return values, if possible
	var returnDataPtr uint32
	malloc := m.ExportedFunction("malloc")
	if malloc != nil {
		res, err := malloc.Call(ctx, 256)
		if err == nil {
			returnDataPtr = uint32(res[0])
		}
	}

	return Module{
		Module:        m,
		ReturnDataPtr: returnDataPtr,
	}
}
