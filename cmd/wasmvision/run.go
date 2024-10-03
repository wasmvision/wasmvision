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
	modelsDir := cCtx.String("models-dir")
	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}
	logging := cCtx.Bool("logging")

	ctx := context.Background()

	// load wasm runtime
	r := runtime.New(ctx, runtime.InterpreterConfig{ModelsDir: modelsDir, Logging: logging})
	defer r.Close(ctx)

	for _, p := range processors {
		module, err := os.ReadFile(p)
		if err != nil {
			log.Panicf("failed to read wasm processor module: %v\n", err)
		}

		if logging {
			log.Printf("Loading wasmCV guest module %s...\n", p)
		}

		if err := r.RegisterGuestModule(ctx, module); err != nil {
			log.Panicf("failed to load wasm processor module: %v\n", err)
		}
	}

	// Open the webcam.
	webcam := capture.NewWebcam(source)
	if err := webcam.Open(); err != nil {
		log.Panicf("Error opening video capture %v\n", source)
	}
	defer webcam.Close()

	switch output {
	case "mjpeg":
		if dest == "" {
			dest = ":8080"
		}
		mjpegstream = engine.NewMJPEGStream(dest)

		go mjpegstream.Start()
	case "file":
		if dest == "" {
			log.Panicf("you must profile a file destination for output-kind=file\n")
		}
		videoWriter = engine.NewVideoWriter(dest)

		if err := videoWriter.Start(webcam); err != nil {
			log.Panicf("Error starting video writer device: %v\n", err)
			return err
		}

		defer videoWriter.Close()
	default:
		log.Panicf("Unknown output kind %v\n", output)
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
			if err := videoWriter.Write(frame); err != nil {
				log.Printf("error writing frame: %v\n", err)
			}
		}

		// cleanup frame
		frame.Close()
		r.FrameCache.Delete(frame.ID)
	}
}
