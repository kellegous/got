package main

import (
	"os"

	"github.com/kellegous/got/pkg/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		os.Exit(1)
	}
}
