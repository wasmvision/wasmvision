package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/runtime"
)

func run(cCtx *cli.Context) error {
	processors := cCtx.StringSlice("processor")
	if len(processors) == 0 {
		fmt.Println("No wasm processors specified")
		os.Exit(1)
	}

	source := cCtx.String("source")
	output := cCtx.String("output")
	dest := cCtx.String("destination")
	processorsDir := cCtx.String("processors-dir")
	if processorsDir == "" {
		processorsDir = DefaultProcessorsPath()
	}
	modelsDir := cCtx.String("models-dir")
	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}
	logging := cCtx.Bool("logging")

	settings := map[string]string{}
	config := cCtx.StringSlice("config")
	for _, c := range config {
		parts := strings.Split(c, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid config format: %v", c)
		}
		settings[parts[0]] = parts[1]
	}

	ctx := context.Background()

	// load wasm runtime
	r := runtime.New(ctx, runtime.InterpreterConfig{
		ProcessorsDir: processorsDir,
		ModelsDir:     modelsDir,
		Logging:       logging,
		Settings:      settings,
	})
	defer r.Close(ctx)

	// load wasm processors
	if err := r.LoadProcessors(ctx, processors); err != nil {
		return fmt.Errorf("failed to load processors: %w", err)
	}

	// Open the webcam.
	webcam := capture.NewWebcam(source)
	if err := webcam.Open(); err != nil {
		return fmt.Errorf("failed opening video capture: %w", err)
	}
	defer webcam.Close()

	var (
		mjpegstream *engine.MJPEGStream
		videoWriter *engine.VideoWriter
	)

	switch output {
	case "mjpeg":
		if dest == "" {
			dest = ":8080"
		}
		mjpegstream = engine.NewMJPEGStream(r.Refs, dest)

		go mjpegstream.Start()
		defer mjpegstream.Close()

	case "file":
		if dest == "" {
			return fmt.Errorf("you must profile a file destination for output=file")
		}
		videoWriter = engine.NewVideoWriter(r.Refs, dest)

		if err := videoWriter.Start(webcam); err != nil {
			return fmt.Errorf("failed starting video writer: %w", err)
		}

		defer videoWriter.Close()
	default:
		return fmt.Errorf("unknown output kind %v", output)
	}

	if logging {
		log.Printf("Reading video from source '%v'\n", source)
	}
	i := 0

	for {
		frame, err := webcam.Read()
		if err != nil {
			switch err {
			case capture.ErrClosed:
				frame.Close()
				if logging {
					log.Println("Capture closed.")
				}

				return nil

			default:
				return fmt.Errorf("frame read error: %w", err)
			}
		}

		if frame.Empty() {
			frame.Close()
			continue
		}

		r.Refs.Set(frame.ID.Unwrap(), frame)

		i++
		if logging {
			log.Printf("Read frame %d\n", i)
		}

		frame = r.Process(ctx, frame)

		switch output {
		case "mjpeg":
			mjpegstream.Publish(frame)
		case "file":
			videoWriter.Write(frame)
		}
	}
}
