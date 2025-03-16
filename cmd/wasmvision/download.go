package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/wasmvision/wasmvision/guest"
	"github.com/wasmvision/wasmvision/models"
)

func downloadModel(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("model name required")
	}
	name := cmd.Args().Get(0)

	if !models.ModelWellKnown(name) {
		return fmt.Errorf("unknown model %s", name)
	}

	modelsDir := cmd.String("models-dir")
	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}

	fmt.Printf("Downloading model %s...\n", name)

	err := models.Download(name, modelsDir)
	if err != nil {
		fmt.Printf("Error downloading model: %s", err)
		return err
	}

	fmt.Printf("Model download complete for %s\n", name)

	return nil
}

func downloadProcessor(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() < 1 {
		return fmt.Errorf("processor name required")
	}
	name := cmd.Args().Get(0)

	if !guest.ProcessorWellKnown(name) {
		return fmt.Errorf("unknown processor %s", name)
	}

	processorsDir := cmd.String("processors-dir")
	if processorsDir == "" {
		processorsDir = DefaultProcessorsPath()
	}

	fmt.Printf("Downloading processor %s...\n", name)

	err := guest.DownloadProcessor(name, processorsDir)
	if err != nil {
		fmt.Printf("Error downloading processor: %s", err)
		return err
	}

	fmt.Printf("Processor download complete for %s\n", name)

	return nil
}
