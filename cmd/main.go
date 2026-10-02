package main

import (
	"fmt"
	"os"

	"fintual/pkg/cli"
)

func main() {
	app := cli.NewCLIApp("data")
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}
}
