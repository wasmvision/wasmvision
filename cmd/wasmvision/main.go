package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wasmvision/wasmvision/capture"
	"github.com/wasmvision/wasmvision/engine"
	"github.com/wasmvision/wasmvision/runtime"
)

var (
	device     = flag.String("device", "/dev/video0", "video capture device to use")
	processors = flag.String("processors", "", "wasm modules to use for processing frames")
	mjpeg      = flag.Bool("mjpeg", false, "output MJPEG stream")
	mjpegPort  = flag.String("mjpeg-port", ":8080", "MJPEG stream port")
)

func main() {
	flag.Parse()

	if *processors == "" {
		log.Panic("processor module is required")
	}

	ctx := context.Background()

	// load wasm runtime
	r := runtime.New(ctx)
	defer r.Close(ctx)

	procs := strings.Split(*processors, ",")
	for _, p := range procs {
		module, err := os.ReadFile(p)
		if err != nil {
			log.Panicf("failed to read wasm processor module: %v\n", err)
		}

		fmt.Printf("Loading wasmCV guest module %s...\n", p)
		r.RegisterGuestModule(ctx, module)
	}

	// Open the webcam.
	webcam := capture.NewWebcam(*device)
	defer webcam.Close()
	if err := webcam.Open(); err != nil {
		log.Panicf("Error opening video capture device: %v\n", *device)
	}

	var mjpegstream engine.MJPEGStream
	if *mjpeg {
		mjpegstream = engine.NewMJPEGStream(*mjpegPort)

		go mjpegstream.Start()
	}

	fmt.Printf("Start reading device: %v\n", *device)
	i := 0
	for {
		frame, err := webcam.Read()
		if err != nil {
			fmt.Printf("frame error %v\n", *device)
			frame.Close()
			continue
		}

		if frame.Image.Empty() {
			frame.Close()
			continue
		}

		engine.FrameCache[frame.ID] = frame

		// clear screen
		fmt.Print("\033[H\033[2J")

		i++
		fmt.Printf("Read frame %d\n", i+1)

		frame = r.PerformProcessing(ctx, frame)

		if *mjpeg {
			mjpegstream.Publish(frame)
		}

		// cleanup frame
		frame.Close()
		delete(engine.FrameCache, frame.ID)
	}
}
