//go:build cuda

package main

import (
	"fmt"

	"gocv.io/x/gocv/cuda"
)

func printDNNBackends() {
	fmt.Print("CV DNN backends:  ")
	fmt.Print("CPU")
	fmt.Println(" GPU")
	devices := cuda.GetCudaEnabledDeviceCount()
	switch devices {
	case 0:
		fmt.Print("  No CUDA devices found")
	default:
		for i := 0; i < devices; i++ {
			fmt.Print("  ")
			cuda.PrintShortCudaDeviceInfo(i)
		}
	}
}
