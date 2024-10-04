package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/runtime"
)

var (
	mjpegstream engine.MJPEGStream
	videoWriter engine.VideoWriter
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

	ctx := context.Background()

	// load wasm runtime
	r := runtime.New(ctx, runtime.InterpreterConfig{
		ProcessorsDir: processorsDir,
		ModelsDir:     modelsDir,
		Logging:       logging,
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

	switch output {
	case "mjpeg":
		if dest == "" {
			dest = ":8080"
		}
		mjpegstream = engine.NewMJPEGStream(r.FrameCache, dest)

		go mjpegstream.Start()
	case "file":
		if dest == "" {
			return fmt.Errorf("you must profile a file destination for output=file")
		}
		videoWriter = engine.NewVideoWriter(r.FrameCache, dest)

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
				return nil

			default:
				fmt.Printf("frame error %v\n", err)
				frame.Close()
				return err
			}
		}

		if frame.Empty() {
			frame.Close()
			continue
		}

		r.FrameCache.Set(frame)

		i++
		if logging {
			log.Printf("Read frame %d\n", i+1)
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
