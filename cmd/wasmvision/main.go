package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

var (
	runFlags = []cli.Flag{
		&cli.StringFlag{Name: "source", Aliases: []string{"s"}, Value: "0", Usage: "video capture source to use such as webcam or file (0 is the default webcam on most systems)"},
		&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Value: "mjpeg", Usage: "output type (mjpeg, file)"},
		&cli.StringFlag{Name: "destination", Aliases: []string{"d"}, Usage: "output destination (port, file path)"},
		&cli.StringSliceFlag{
			Name:    "processor",
			Aliases: []string{"p"},
			Usage:   "wasm module to use for processing frames. Format: -processor /path/processor1.wasm -processor /path2/processor2.wasm",
		},
		&cli.BoolFlag{Name: "logging", Value: true, Usage: "log detailed info to console (default: true)"},
		&cli.StringFlag{Name: "models-dir", Aliases: []string{"models"}, EnvVars: []string{"WASMVISION_MODELS_DIR"}, Usage: "directory for model loading (default to $home/models)"},
		&cli.BoolFlag{Name: "model-download", Aliases: []string{"download"}, Value: true, Usage: "automatically download known models (default: true)"},
		&cli.StringFlag{Name: "processors-dir", Aliases: []string{"processors"}, EnvVars: []string{"WASMVISION_PROCESSORS_DIR"}, Usage: "directory for processor loading (default to $home/processors)"},
		&cli.BoolFlag{Name: "processor-download", Value: true, Usage: "automatically download known processors (default: true)"},
	}

	downloadFlags = []cli.Flag{
		&cli.StringFlag{Name: "models-dir", Aliases: []string{"models"}, EnvVars: []string{"WASMVISION_MODELS_DIR"}, Usage: "directory for model loading (default to $home/models)"},
		&cli.BoolFlag{Name: "processor-download", Value: true, Usage: "automatically download known processors (default: true)"},
	}
)

func main() {
	app := &cli.App{
		Name:    "wasmvision",
		Usage:   "wasmVision CLI",
		Version: Version(),
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
				Subcommands: []*cli.Command{
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
				Subcommands: []*cli.Command{
					{
						Name:   "models",
						Usage:  "lists all known computer vision models",
						Action: listallModels,
					},
					{
						Name:   "processors",
						Usage:  "lists all known wasm processors",
						Action: listallProcessors,
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

	if err := app.Run(os.Args); err != nil {
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
