package main

import (
	"os"

	"github.com/omegatymbjiep/ilab1/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
