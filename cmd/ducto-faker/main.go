// cmd/ducto-faker/main.go
package main

import (
	"github.com/tommed/ducto-faker/internal/cli"
	"os"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
