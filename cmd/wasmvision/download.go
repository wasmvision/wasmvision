package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/wasmvision/wasmvision/net"
)

func download(cCtx *cli.Context) error {
	if cCtx.Args().Len() < 1 {
		return fmt.Errorf("model name required")
	}
	name := cCtx.Args().Get(0)

	dl, ok := net.KnownModels[name]
	if !ok {
		return fmt.Errorf("unknown model %s", name)
	}

	modelsDir := cCtx.String("models-dir")
	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}

	fmt.Printf("Downloading model %s...\n", name)

	err := net.DownloadModel(dl, modelsDir)
	if err != nil {
		fmt.Printf("Error downloading model: %s", err)
		return err
	}

	fmt.Printf("Model download complete for %s\n", name)

	return nil
}
