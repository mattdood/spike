package main

import (
	"os"

	"github.com/mattdood/spike/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
