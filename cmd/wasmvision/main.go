package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hybridgroup/wasmvision/capture"
	"github.com/hybridgroup/wasmvision/engine"
	"github.com/hybridgroup/wasmvision/runtime"
)

var (
	device     = flag.String("device", "/dev/video0", "video capture device to use")
	moduleName = flag.String("module", "", "wasm module to use for processing frames")
)

func main() {
	flag.Parse()

	if *moduleName == "" {
		log.Panic("processor module is required")
	}

	module, err := os.ReadFile(*moduleName)
	if err != nil {
		log.Panicf("failed to read wasm processor module: %v\n", err)
	}

	ctx := context.Background()

	// load wasm runtime
	r := runtime.New(ctx)
	defer r.Close(ctx)

	fmt.Printf("Loading wasmCV guest module %s...\n", *moduleName)
	mod, err := r.Instantiate(ctx, module)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}
	process := mod.ExportedFunction("process")

	// Open the webcam.
	webcam := capture.NewWebCam(*device)
	defer webcam.Close()
	if err := webcam.Open(); err != nil {
		log.Panicf("Error opening video capture device: %v\n", *device)
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
		fmt.Printf("Running wasmCV module %s\n", *moduleName)

		if _, err := process.Call(ctx, uint64(frame.ID)); err != nil {
			fmt.Printf("Error calling process: %v\n", err)
		}

		// cleanup frame
		frame.Close()
		delete(engine.FrameCache, frame.ID)
	}
}
