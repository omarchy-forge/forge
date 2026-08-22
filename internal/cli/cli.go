// Package cli implements the omaforge command-line interface.
package cli

import (
	"fmt"
	"io"
)

const helpText = `Omarchy Forge
Developer tools for building, testing, and shipping Omarchy plugins.

Usage:
  omaforge [--help]
  omaforge version

Commands:
  version    Print build and compatibility information

Options:
  -h, --help Print this help

Omarchy Forge is in early development. Plugin scaffolding and checking are not
available yet.
`

// BuildInfo contains release metadata supplied at link time.
type BuildInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdout, stderr io.Writer, info BuildInfo) int {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "--help" || args[0] == "-h")) {
		fmt.Fprint(stdout, helpText)
		return 0
	}

	switch args[0] {
	case "version":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "omaforge version does not accept arguments")
			fmt.Fprintln(stderr, "Run 'omaforge --help' for usage.")
			return 2
		}
		printVersion(stdout, info)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		fmt.Fprintln(stderr, "Run 'omaforge --help' for usage.")
		return 2
	}
}

func printVersion(w io.Writer, info BuildInfo) {
	version := valueOrUnknown(info.Version)
	if info.Version == "" {
		version = "dev"
	}
	fmt.Fprintf(w, "omaforge %s\n", version)
	fmt.Fprintf(w, "commit: %s\n", valueOrUnknown(info.Commit))
	fmt.Fprintf(w, "built: %s\n", valueOrUnknown(info.BuildDate))
	fmt.Fprintln(w, "manifest schemas: 1")
}

func valueOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}
