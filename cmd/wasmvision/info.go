package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/urfave/cli/v3"
	"gocv.io/x/gocv"
)

func info(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("wasmVision version %s %s/%s\n", Version(), runtime.GOOS, runtime.GOARCH)

	fmt.Print("Camera backends: ")
	printCameraBackends()
	fmt.Println()

	fmt.Print("Stream backends: ")
	printStreamBackends()
	fmt.Println()

	fmt.Print("Writer backends: ")
	printWriterBackends()
	fmt.Println()

	return nil
}

func printCameraBackends() {
	cameraBacks := gocv.VideoRegistry.GetCameraBackends()
	for _, b := range cameraBacks {
		fmt.Printf(" %s", gocv.VideoRegistry.GetBackendName(b))
	}
}

func printStreamBackends() {
	streamBacks := gocv.VideoRegistry.GetStreamBackends()
	for _, b := range streamBacks {
		fmt.Printf(" %s", gocv.VideoRegistry.GetBackendName(b))
	}
}

func printWriterBackends() {
	writerBacks := gocv.VideoRegistry.GetWriterBackends()
	for _, b := range writerBacks {
		fmt.Printf(" %s", gocv.VideoRegistry.GetBackendName(b))
	}
}
