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
	}

	downloadFlags = []cli.Flag{
		&cli.StringFlag{Name: "models-dir", Aliases: []string{"models"}, EnvVars: []string{"WASMVISION_MODELS_DIR"}, Usage: "directory for model loading (default to $home/models)"},
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
				Name:      "download",
				Usage:     "Download computer vision models",
				ArgsUsage: "[known-model-name]",
				Action:    download,
				Flags:     downloadFlags,
			},
			{
				Name:   "info",
				Usage:  "Show installation info",
				Action: info,
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

func DefaultModelPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dirname, "models")
}
