package main

import (
"fmt"
"os"

"github.com/ecoker/launchpad/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}
}
