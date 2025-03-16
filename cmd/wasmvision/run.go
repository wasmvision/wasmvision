package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/runtime"
)

func run(ctx context.Context, cmd *cli.Command) error {
	processors := cmd.StringSlice("processor")
	if len(processors) == 0 {
		fmt.Println("No wasm processors specified")
		os.Exit(1)
	}

	source := cmd.String("source")
	output := cmd.String("output")
	dest := cmd.String("destination")
	processorsDir := cmd.String("processors-dir")
	if processorsDir == "" {
		processorsDir = DefaultProcessorsPath()
	}
	modelsDir := cmd.String("models-dir")
	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}

	logging := cmd.String("logging")
	switch logging {
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		return fmt.Errorf("unknown log level %v", logging)
	}

	settings := map[string]string{}
	config := cmd.StringSlice("config")
	for _, c := range config {
		parts := strings.Split(c, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid config format: %v", c)
		}
		settings[parts[0]] = parts[1]
	}

	// load wasm runtime
	r := runtime.New(ctx, runtime.InterpreterConfig{
		ProcessorsDir: processorsDir,
		ModelsDir:     modelsDir,
		Settings:      settings,
	})
	defer r.Close(ctx)

	// load wasm processors
	if err := r.LoadProcessors(ctx, processors); err != nil {
		return fmt.Errorf("failed to load processors: %w", err)
	}

	// Open the capture device.
	cap := cmd.String("capture")
	var device capture.Capture

	switch cap {
	case "auto", "webcam":
		device = capture.NewWebcam(source)
		if err := device.Open(); err != nil {
			return fmt.Errorf("failed opening video capture: %w", err)
		}
	case "gstreamer":
		device = capture.NewGStreamer(source)
		if err := device.Open(); err != nil {
			return fmt.Errorf("failed opening video capture stream: %w", err)
		}
	case "ffmpeg":
		device = capture.NewFFmpeg(source)
		if err := device.Open(); err != nil {
			return fmt.Errorf("failed opening video capture stream: %w", err)
		}
	default:
		return fmt.Errorf("unknown capture type %v", cap)
	}

	defer device.Close()

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

		if err := mjpegstream.Start(); err != nil {
			return fmt.Errorf("failed starting mjpeg stream: %w", err)
		}

		defer mjpegstream.Close()

	case "file":
		if dest == "" {
			return fmt.Errorf("you must profile a file destination for output=file")
		}
		videoWriter = engine.NewVideoWriter(r.Refs, dest)

		if err := videoWriter.Start(device); err != nil {
			return fmt.Errorf("failed starting video writer: %w", err)
		}

		defer videoWriter.Close()
	default:
		return fmt.Errorf("unknown output kind %v", output)
	}

	slog.Info(fmt.Sprintf("Reading video from source '%v", source))
	i := 0

	for {
		frame, err := device.Read()
		if err != nil {
			switch err {
			case capture.ErrClosed:
				frame.Close()
				slog.Info("Capture closed.")

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
		slog.Info(fmt.Sprintf("Read frame %d", i))

		frame = r.Process(ctx, frame)

		switch output {
		case "mjpeg":
			mjpegstream.Publish(frame)
		case "file":
			videoWriter.Write(frame)
		}
	}
}
