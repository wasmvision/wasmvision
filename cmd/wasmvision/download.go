package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
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

	return downloadModel(dl, modelsDir)
}

func downloadModel(model net.ModelFile, targetDir string) error {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting pwd: %s", err)
		return err
	}

	opts := []getter.ClientOption{}
	client := &getter.Client{
		Ctx:     context.Background(),
		Src:     model.URL,
		Dst:     filepath.Join(targetDir, filepath.Base(model.Filename)),
		Pwd:     pwd,
		Mode:    getter.ClientModeFile,
		Options: opts,
	}

	if err := client.Get(); err != nil {
		fmt.Printf("Error downloading model: %s", err)
		return err
	}

	fmt.Printf("Model download complete for %s\n", model.Filename)

	return nil
}
