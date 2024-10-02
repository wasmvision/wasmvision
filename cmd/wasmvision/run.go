package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/runtime"
	"gocv.io/x/gocv"
)

var (
	mjpegstream engine.MJPEGStream
	videoWriter *gocv.VideoWriter
)

func run(cCtx *cli.Context) error {
	processors := cCtx.StringSlice("processor")
	if len(processors) == 0 {
		fmt.Println("No wasm processors specified")
		os.Exit(1)
	}

	source := cCtx.String("source")
	output := cCtx.String("output-kind")
	dest := cCtx.String("destination")
	clear := cCtx.Bool("clear")
	modelsDir := cCtx.String("models-dir")
	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}

	ctx := context.Background()

	// load wasm runtime
	r := runtime.New(ctx, runtime.InterpreterConfig{ModelsDir: modelsDir})
	defer r.Close(ctx)

	for _, p := range processors {
		module, err := os.ReadFile(p)
		if err != nil {
			log.Panicf("failed to read wasm processor module: %v\n", err)
		}

		fmt.Printf("Loading wasmCV guest module %s...\n", p)
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
		var err error
		videoWriter, err = openVideoWriter(webcam, dest)
		if err != nil {
			log.Panicf("Error opening video file writer device: %v\n", err)
			return err
		}
		defer videoWriter.Close()
	default:
		fmt.Printf("Unknown output kind %v\n", output)
	}

	fmt.Printf("Reading video from %v\n", source)
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

		if clear {
			fmt.Print("\033[2J\033[3J\033[H")
		}

		i++
		fmt.Printf("Read frame %d\n", i+1)

		frame = r.Process(ctx, frame)

		switch output {
		case "mjpeg":
			mjpegstream.Publish(frame)
		case "file":
			if err := videoWriter.Write(frame.Image); err != nil {
				fmt.Printf("error writing frame: %v\n", err)
			}
		}

		// cleanup frame
		frame.Close()
		r.FrameCache.Delete(frame.ID)
	}
}

func openVideoWriter(source capture.Capture, dest string) (*gocv.VideoWriter, error) {
	sample, err := source.Read()
	if err != nil {
		fmt.Printf("frame error %v\n", err)
		return nil, err
	}

	defer sample.Close()

	videoWriter, err := gocv.VideoWriterFile(dest, "MJPG", 25, sample.Image.Cols(), sample.Image.Rows(), true)
	if err != nil {
		fmt.Printf("error opening video file writer device: %v\n", err)
		return nil, err
	}

	return videoWriter, nil
}

func DefaultModelPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(dirname, "models")
}
