package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	runFlags = []cli.Flag{
		&cli.StringFlag{Name: "device", Aliases: []string{"d"}, Value: "0", Usage: "video capture device to use (default: 0)"},
		&cli.BoolFlag{Name: "mjpeg", Usage: "output MJPEG stream (default: false)"},
		&cli.StringFlag{Name: "mjpegport", Usage: "MJPEG stream port (default :8080)", Value: ":8080"},
		&cli.StringSliceFlag{
			Name:    "processor",
			Aliases: []string{"p"},
			Usage:   "wasm module to use for processing frames. Format: -processor /path/processor1.wasm -processor /path2/processor2.wasm",
		},
		&cli.BoolFlag{Name: "clear-screen", Aliases: []string{"clear"}, Value: true, Usage: "clear screen between frames (default: true)"},
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
