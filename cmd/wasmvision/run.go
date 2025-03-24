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
	if len(pipeline) > 0 {
		list := pipeline[0]
		list = strings.TrimLeft(list, "[")
		list = strings.TrimRight(list, "]")
		processors = strings.Split(list, " ")
	}
	if len(processors) == 0 {
		fmt.Println("No wasm processors specified")
		os.Exit(1)
	}

	if processorsDir == "" {
		processorsDir = DefaultProcessorsPath()
	}

	if modelsDir == "" {
		modelsDir = DefaultModelPath()
	}

	switch loggingLevel {
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		return fmt.Errorf("unknown log level %v", loggingLevel)
	}

	if len(configuration) > 0 {
		list := configuration[0]
		list = strings.TrimLeft(list, "[")
		list = strings.TrimRight(list, "]")
		config = strings.Split(list, " ")
	}
	settings := map[string]string{}
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
	var device capture.Capture

	switch captureDevice {
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
		return fmt.Errorf("unknown capture type %v", captureDevice)
	}

	defer device.Close()

	var (
		mjpegstream *engine.MJPEGStream
		videoWriter *engine.VideoWriter
	)

	switch output {
	case "mjpeg":
		if destination == "" {
			destination = ":8080"
		}
		mjpegstream = engine.NewMJPEGStream(r.Refs, destination)

		if err := mjpegstream.Start(); err != nil {
			return fmt.Errorf("failed starting mjpeg stream: %w", err)
		}

		defer mjpegstream.Close()

	case "file":
		if destination == "" {
			return fmt.Errorf("you must profile a file destination for output=file")
		}
		videoWriter = engine.NewVideoWriter(r.Refs, destination)

		if err := videoWriter.Start(device); err != nil {
			return fmt.Errorf("failed starting video writer: %w", err)
		}

		defer videoWriter.Close()
	default:
		return fmt.Errorf("unknown output kind %v", output)
	}

	var (
		mcpServer *engine.MCPServer
	)
	if mcpEnabled {
		mcpServer = engine.NewMCPServer(mcpPort)
		if err := mcpServer.Start(); err != nil {
			return fmt.Errorf("failed starting MCP server: %w", err)
		}

		defer mcpServer.Close()
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

		if mcpEnabled {
			if err := mcpServer.PublishInput(frame); err != nil {
				slog.Error("failed to publish input frame:" + err.Error())
			}
		}

		outframe := r.Process(ctx, frame)

		if mcpEnabled {
			if err := mcpServer.PublishOutput(outframe); err != nil {
				slog.Error("failed to publish output frame:" + err.Error())
			}
		}

		switch output {
		case "mjpeg":
			mjpegstream.Publish(outframe)
		case "file":
			videoWriter.Write(outframe)
		}
		frame.Close()
	}
}
