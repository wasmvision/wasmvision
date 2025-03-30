//go:build !cuda

package main

import (
	"fmt"
)

func printDNNBackends() {
	fmt.Print("CV DNN backends:  ")
	fmt.Println("CPU")
}
