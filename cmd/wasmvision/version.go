package main

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
)

var (
	v   = "0.1.0-beta3"
	sha string
)

func version(cCtx *cli.Context) error {
	fmt.Printf("wasmVision version %s %s/%s\n", Version(), runtime.GOOS, runtime.GOARCH)

	return nil
}

func Version() string {
	if sha != "" {
		return v + "-" + sha
	}
	return v
}
