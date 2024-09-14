package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hybridgroup/wasmvision/engine"
	"github.com/tetratelabs/wazero"
	"gocv.io/x/gocv"
)

var (
	device    = flag.String("device", "/dev/video0", "video capture device to use")
	processor = flag.String("processor", "", "wasm module to use for processing frames")
)

func main() {
	flag.Parse()

	if *processor == "" {
		log.Panic("processor is required")
	}

	module, err := os.ReadFile(*processor)
	if err != nil {
		log.Panicf("failed to read wasm module: %v\n", err)
	}

	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	println("Defining host functions...")
	modules := engine.HostModules()
	if err := modules.DefineWazero(r, nil); err != nil {
		log.Panicf("error define host functions: %v\n", err)
	}

	fmt.Printf("Loading wasmCV guest module %s...\n", *processor)
	mod, err := r.Instantiate(ctx, module)
	if err != nil {
		log.Panicf("failed to instantiate module: %v", err)
	}
	process := mod.ExportedFunction("process")

	// Open the webcam.
	webcam, err := gocv.OpenVideoCapture(*device)
	if err != nil {
		log.Panicf("Error opening video capture device: %v\n", *device)
	}
	defer webcam.Close()

	fmt.Printf("Start reading device: %v\n", *device)
	i := 0
	for {
		img := gocv.NewMat()

		if ok := webcam.Read(&img); !ok {
			fmt.Printf("frame error %v\n", *device)
			continue
		}

		if img.Empty() {
			continue
		}

		frame := engine.NewFrame()
		frame.SetImage(img)
		engine.FrameCache[frame.ID] = frame

		// clear screen
		fmt.Print("\033[H\033[2J")

		i++
		fmt.Printf("Read frame %d\n", i+1)
		fmt.Printf("Running wasmCV module %s\n", *processor)

		_, err := process.Call(ctx, uint64(frame.ID))
		if err != nil {
			fmt.Printf("Error calling process: %v\n", err)
		}

		// cleanup frame
		frame.Close()
		delete(engine.FrameCache, frame.ID)
	}
}
