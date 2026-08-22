package main

import (
	"os"

	"github.com/omarchy-forge/forge/internal/cli"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	info := cli.BuildInfo{
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr, info))
}
