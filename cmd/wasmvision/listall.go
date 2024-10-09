package main

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/wasmvision/wasmvision/guest"
	"github.com/wasmvision/wasmvision/models"
)

func listallModels(cCtx *cli.Context) error {
	printModels()
	fmt.Println()

	return nil
}

func listallProcessors(cCtx *cli.Context) error {
	printProcessors()
	fmt.Println()

	return nil
}

type keyValue struct {
	key   string
	value string
}

func printModels() {
	s := make([]keyValue, 0, len(models.KnownModels))
	for k, v := range models.KnownModels {
		s = append(s, keyValue{k, v.Alias})
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].value < s[j].value
	})

	// iterate over the slice to get the desired order
	for _, v := range s {
		fmt.Println(v.value)
	}
}

func printProcessors() {
	s := make([]keyValue, 0, len(guest.KnownProcessors))
	for k, v := range guest.KnownProcessors {
		s = append(s, keyValue{k, v.Alias})
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].value < s[j].value
	})

	// iterate over the slice to get the desired order
	for _, v := range s {
		fmt.Println(v.value)
	}
}
