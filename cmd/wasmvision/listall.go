package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/urfave/cli/v3"
	"github.com/wasmvision/wasmvision/guest"
	"github.com/wasmvision/wasmvision/models"
)

func listallModels(ctx context.Context, cmd *cli.Command) error {
	printModels(cmd.Bool("long"))
	fmt.Println()

	return nil
}

func listallProcessors(ctx context.Context, cmd *cli.Command) error {
	printProcessors(cmd.Bool("long"))
	fmt.Println()

	return nil
}

type modelValue struct {
	key   string
	value models.ModelFile
}

func printModels(long bool) {
	s := make([]modelValue, 0, len(models.KnownModels))
	for k, v := range models.KnownModels {
		s = append(s, modelValue{k, v})
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].value.Alias < s[j].value.Alias
	})

	for _, v := range s {
		if long {
			fmt.Printf("%-30s  %s\n", v.value.Alias, v.value.Description)
		} else {
			fmt.Println(v.value.Alias)
		}
	}
}

type processorValue struct {
	key   string
	value guest.ProcessorFile
}

func printProcessors(long bool) {
	s := make([]processorValue, 0, len(guest.KnownProcessors()))
	for k, v := range guest.KnownProcessors() {
		s = append(s, processorValue{k, v})
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].value.Alias < s[j].value.Alias
	})

	// iterate over the slice to get the desired order
	for _, v := range s {
		if long {
			fmt.Printf("%-20s  %s\n", v.value.Alias, v.value.Description)
		} else {
			fmt.Println(v.value.Alias)
		}
	}
}
