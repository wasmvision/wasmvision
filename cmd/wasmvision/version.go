package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/urfave/cli/v3"
)

var (
	v   = "0.4.0-dev"
	sha string
)

func version(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("wasmVision version %s %s/%s\n", Version(), runtime.GOOS, runtime.GOARCH)

	return nil
}

func Version() string {
	if sha != "" {
		return v + "-" + sha
	}
	return v
}
