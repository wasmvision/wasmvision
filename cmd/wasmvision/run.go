package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/runtime"
)

var (
	mjpegstream *engine.MJPEGStream
	videoWriter *engine.VideoWriter
)

func run(ctx context.Context, cmd *cli.Command) error {
	handlePipelineParams()
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

	if err := setLoggingLevel(); err != nil {
		return err
	}

	if err := handleConfigurationParams(); err != nil {
		return err
	}

	if enableCUDA {
		if !runtime.CheckCUDA() {
			return fmt.Errorf("CUDA not available on this system")
		}
		slog.Info("CUDA enabled")
	}

	// load wasm runtime
	r, err := runtime.New(ctx, runtime.InterpreterConfig{
		ProcessorsDir: processorsDir,
		ModelsDir:     modelsDir,
		Settings:      config,
		Datastorage:   datastorage,
		EnableCUDA:    enableCUDA,
	})
	if err != nil {
		return fmt.Errorf("failed to create runtime: %w", err)
	}
	defer r.Close(ctx)

	// load wasm processors
	if err := r.LoadProcessors(ctx, processors); err != nil {
		return fmt.Errorf("failed to load processors: %w", err)
	}

	// Open the capture device.
	device, err := capture.OpenDevice(captureDevice, source)
	if err != nil {
		return fmt.Errorf("failed to open capture device: %w", err)
	}
	defer device.Close()

	if height > 0 {
		if err := device.Set(capture.VideoCaptureFrameHeight, float32(height)); err != nil {
			return fmt.Errorf("failed to set height: %w", err)
		}
	}
	if width > 0 {
		if err := device.Set(capture.VideoCaptureFrameWidth, float32(width)); err != nil {
			return fmt.Errorf("failed to set width: %w", err)
		}
	}

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
	case "none":
		// do nothing
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
				slog.Error("failed to publish input frame to MCP server:" + err.Error())
			}
		}

		outframe, err := r.Process(ctx, frame)
		if err != nil {
			slog.Error("failed to process frame: " + err.Error())
			frame.Close()

			return err
		}

		if mcpEnabled {
			if err := mcpServer.PublishOutput(outframe); err != nil {
				slog.Error("failed to publish output frame to MCP server:" + err.Error())
			}
		}

		switch output {
		case "mjpeg":
			mjpegstream.Publish(outframe)
		case "file":
			videoWriter.Write(outframe)
		case "none":
			outframe.Close()
		}

		// Close the original frame unless it was returned by the output
		if frame.ID.Unwrap() != outframe.ID.Unwrap() {
			frame.Close()
		}
	}
}
