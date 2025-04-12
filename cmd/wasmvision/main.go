package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:        "wasmvision",
		Usage:       "wasmVision CLI",
		Description: "wasmVision gets you up and running with computer vision.",
		Version:     Version(),
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "Run wasmVision processors",
				Action: run,
				Flags:  runFlags,
			},
			{
				Name:  "download",
				Usage: "Download computer vision models and processors",
				Commands: []*cli.Command{
					{
						Name:      "model",
						Usage:     "download a known computer vision model",
						ArgsUsage: "[known-model-name]",
						Action:    downloadModel,
						Flags:     downloadFlags,
					},
					{
						Name:      "processor",
						Usage:     "download a known processor",
						ArgsUsage: "[known-processor-name]",
						Action:    downloadProcessor,
						Flags:     downloadFlags,
					},
				},
			},
			{
				Name:   "info",
				Usage:  "Show installation info",
				Action: info,
			},
			{
				Name:  "listall",
				Usage: "Lists all known models and processors",
				Commands: []*cli.Command{
					{
						Name:   "models",
						Usage:  "lists all known computer vision models",
						Action: listallModels,
						Flags:  listAllFlags,
					},
					{
						Name:   "processors",
						Usage:  "lists all known wasm processors",
						Action: listallProcessors,
						Flags:  listAllFlags,
					},
				},
			},
			{
				Name:   "version",
				Usage:  "Show version",
				Action: version,
			},
			{
				Name:   "about",
				Usage:  "About wasmVision",
				Action: about,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func DefaultProcessorsPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dirname, "processors")
}

func DefaultModelPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dirname, "models")
}
