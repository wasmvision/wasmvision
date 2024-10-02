package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	runFlags = []cli.Flag{
		&cli.StringFlag{Name: "source", Aliases: []string{"s"}, Value: "0", Usage: "video capture source to use such as webcam or file (0 is the default webcam on most systems)"},
		&cli.StringFlag{Name: "output-kind", Aliases: []string{"o"}, Usage: "kind of output (mjpeg, file)"},
		&cli.StringFlag{Name: "destination", Aliases: []string{"d"}, Usage: "destination for the output (port, file path)"},
		&cli.StringSliceFlag{
			Name:    "processor",
			Aliases: []string{"p"},
			Usage:   "wasm module to use for processing frames. Format: -processor /path/processor1.wasm -processor /path2/processor2.wasm",
		},
		&cli.BoolFlag{Name: "clear-screen", Aliases: []string{"clear"}, Value: true, Usage: "clear screen between frames (default: true)"},
		&cli.StringFlag{Name: "models-dir", Aliases: []string{"models"}, EnvVars: []string{"WASMVISION_MODELS_DIR"}, Usage: "Directory for model loading (default to $home/models)"},
		&cli.BoolFlag{Name: "model-download", Aliases: []string{"download"}, Value: true, Usage: "automatically download known models (default: true)"},
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
				Name:   "download",
				Usage:  "Download computer vision models",
				Action: download,
				Flags:  runFlags,
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
