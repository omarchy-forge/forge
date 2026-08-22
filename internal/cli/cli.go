// Package cli implements the omaforge command-line interface.
package cli

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/omarchy-forge/forge/internal/checks"
	"github.com/omarchy-forge/forge/internal/dev"
	"github.com/omarchy-forge/forge/internal/doctor"
	"github.com/omarchy-forge/forge/internal/scaffold"
)

const helpText = `Omarchy Forge
Developer tools for building, testing, and shipping Omarchy plugins.

Usage:
  omaforge [--help]
  omaforge version
  omaforge init <directory> [options]
  omaforge check <directory> [options]
  omaforge doctor [directory]
  omaforge dev <directory> --trust-plugin-code [--state ready|empty|error] [--watch]
  omaforge screenshot <directory> --trust-plugin-code --state ready|empty|error --output <file.png>

Commands:
  init       Create a bar-widget plugin project
  check      Run deterministic, non-executing plugin checks
  doctor     Diagnose the local Omarchy environment and plugin
  dev        Run a trusted plugin in its isolated runtime harness
  screenshot Capture a fictional plugin state without capturing the desktop
  version    Print build and compatibility information

Options:
  -h, --help Print this help

Run 'omaforge <command> --help' for command-specific options.
`

// BuildInfo contains release metadata supplied at link time.
type BuildInfo struct {
	Version   string
	Commit    string
	BuildDate string
}

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer, info BuildInfo) int {
	if len(args) == 0 || (len(args) == 1 && (args[0] == "--help" || args[0] == "-h")) {
		fmt.Fprint(stdout, helpText)
		return 0
	}

	switch args[0] {
	case "init":
		return runInit(args[1:], stdin, stdout, stderr)
	case "check":
		return runCheck(args[1:], stdout, stderr)
	case "doctor":
		return runDoctor(args[1:], stdout, stderr)
	case "dev":
		return runDev(args[1:], stdin, stdout, stderr)
	case "screenshot":
		return runScreenshot(args[1:], stdin, stdout, stderr)
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

func runScreenshot(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	for _, argument := range args {
		if argument == "-h" || argument == "--help" {
			fmt.Fprintln(stdout, "Usage: omaforge screenshot <directory> --trust-plugin-code --state ready|empty|error --output <file.png>")
			return 0
		}
	}
	trusted, directory, state, output := false, "", "", ""
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--trust-plugin-code":
			trusted = true
		case "--state", "--output":
			name := args[index]
			if index+1 >= len(args) {
				fmt.Fprintf(stderr, "omaforge screenshot: %s requires a value\n", name)
				return 2
			}
			index++
			if name == "--state" {
				state = args[index]
			} else {
				output = args[index]
			}
		default:
			if strings.HasPrefix(args[index], "-") || directory != "" {
				fmt.Fprintln(stderr, "omaforge screenshot: invalid arguments")
				return 2
			}
			directory = args[index]
		}
	}
	if !trusted || directory == "" || output == "" || (state != "ready" && state != "empty" && state != "error") {
		fmt.Fprintln(stderr, "omaforge screenshot requires a directory, explicit trust, a ready|empty|error state, and PNG output")
		return 2
	}
	if err := dev.Screenshot(directory, state, output, stdin, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "omaforge screenshot: %v\n", err)
		return 1
	}
	fmt.Fprintln(stdout, "Saved plugin-only screenshot to "+output)
	return 0
}

func runDev(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	for _, argument := range args {
		if argument == "-h" || argument == "--help" {
			fmt.Fprintln(stdout, "Usage: omaforge dev <directory> --trust-plugin-code [--state ready|empty|error] [--watch]")
			fmt.Fprintln(stdout, "Runs the project's isolated runtime harness after explicit trust acknowledgement.")
			return 0
		}
	}
	trusted := false
	directory := ""
	state := ""
	watch := false
	for index := 0; index < len(args); index++ {
		argument := args[index]
		switch argument {
		case "--trust-plugin-code":
			if trusted {
				fmt.Fprintln(stderr, "omaforge dev accepts --trust-plugin-code only once")
				return 2
			}
			trusted = true
		case "--state":
			if index+1 >= len(args) || state != "" {
				fmt.Fprintln(stderr, "omaforge dev --state requires one value")
				return 2
			}
			index++
			state = args[index]
		case "--watch":
			if watch {
				fmt.Fprintln(stderr, "omaforge dev accepts --watch only once")
				return 2
			}
			watch = true
		default:
			if strings.HasPrefix(argument, "-") || directory != "" {
				fmt.Fprintln(stderr, "omaforge dev accepts exactly one plugin directory and --trust-plugin-code")
				return 2
			}
			directory = argument
		}
	}
	if directory == "" || !trusted {
		fmt.Fprintln(stderr, "omaforge dev requires exactly one <directory> and --trust-plugin-code")
		fmt.Fprintln(stderr, "Review the plugin QML and local commands before running it.")
		return 2
	}
	if state != "" && state != "ready" && state != "empty" && state != "error" {
		fmt.Fprintf(stderr, "omaforge dev: unknown state %q; available: ready, empty, error\n", state)
		return 2
	}
	if watch {
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		if err := dev.Watch(ctx, directory, state, stdin, stdout, stderr); err != nil {
			fmt.Fprintf(stderr, "omaforge dev: %v\n", err)
			return 1
		}
		return 0
	}
	if err := dev.Run(directory, state, stdin, stdout, stderr); err != nil {
		fmt.Fprintf(stderr, "omaforge dev: %v\n", err)
		return 1
	}
	return 0
}

func runCheck(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	format := flags.String("format", "text", "report format: text, json, or sarif")
	compatibility := flags.String("omarchy-version", "4", "target Omarchy compatibility version")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: omaforge check <directory> [--format text|json|sarif] [--omarchy-version 4]")
	}
	for _, argument := range args {
		if argument == "-h" || argument == "--help" {
			flags.SetOutput(stdout)
			flags.Usage()
			return 0
		}
	}
	ordered, err := reorderCommandArgs(args, map[string]bool{"--format": true, "--omarchy-version": true})
	if err != nil {
		fmt.Fprintf(stderr, "omaforge check: %v\n", err)
		return 2
	}
	if err := flags.Parse(ordered); err != nil {
		return 2
	}
	if flags.NArg() != 1 {
		fmt.Fprintln(stderr, "omaforge check requires exactly one plugin directory")
		return 2
	}
	if *compatibility != "4" {
		fmt.Fprintf(stderr, "unsupported Omarchy compatibility version %q; supported: 4\n", *compatibility)
		return 2
	}
	report := checks.RunFor(flags.Arg(0), "omarchy-"+*compatibility)
	switch *format {
	case "text":
		err = checks.WriteText(stdout, report)
	case "json":
		err = checks.WriteJSON(stdout, report)
	case "sarif":
		err = checks.WriteSARIF(stdout, report)
	default:
		fmt.Fprintf(stderr, "unknown format %q; available: text, json, sarif\n", *format)
		return 2
	}
	if err != nil {
		fmt.Fprintf(stderr, "omaforge check: write report: %v\n", err)
		return 1
	}
	if report.Summary.Errors > 0 {
		return 1
	}
	return 0
}

func runDoctor(args []string, stdout, stderr io.Writer) int {
	for _, argument := range args {
		if argument == "-h" || argument == "--help" {
			fmt.Fprintln(stdout, "Usage: omaforge doctor [directory]")
			return 0
		}
	}
	if len(args) > 1 {
		fmt.Fprintln(stderr, "omaforge doctor accepts at most one plugin directory")
		return 2
	}
	target := "."
	if len(args) == 1 {
		target = args[0]
	}
	report := doctor.Run(target, doctor.ExecRunner{})
	if err := doctor.WriteText(stdout, report); err != nil {
		fmt.Fprintf(stderr, "omaforge doctor: write report: %v\n", err)
		return 1
	}
	if report.Failed() {
		return 1
	}
	return 0
}

func runInit(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var options scaffold.Options
	var templateName string
	var nonInteractive bool
	var noCI bool
	flags.StringVar(&options.ID, "id", "", "namespaced plugin ID (required)")
	flags.StringVar(&options.Name, "name", "", "plugin display name")
	flags.StringVar(&options.Description, "description", "", "short plugin description")
	flags.StringVar(&options.Author, "author", "", "plugin author (required)")
	flags.StringVar(&options.License, "license", "MIT", "plugin license (MIT only)")
	flags.StringVar(&options.Section, "section", "right", "default bar section: left, center, or right")
	flags.StringVar(&templateName, "template", "bar-widget", "template name (bar-widget only)")
	flags.BoolVar(&options.DryRun, "dry-run", false, "show the write plan without creating files")
	flags.BoolVar(&options.Force, "force", false, "overwrite colliding generated files after showing the plan")
	flags.BoolVar(&options.InitGit, "git", false, "initialize a local Git repository")
	flags.BoolVar(&options.AgentReady, "agent-ready", false, "add a structured specification and agent safety guidance")
	flags.BoolVar(&noCI, "no-ci", false, "omit the generated GitHub Actions workflow")
	flags.BoolVar(&nonInteractive, "non-interactive", false, "fail instead of prompting for missing values")
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: omaforge init <directory> [options]")
		fmt.Fprintln(flags.Output())
		fmt.Fprintln(flags.Output(), "Create a deterministic Omarchy bar-widget plugin project.")
		fmt.Fprintln(flags.Output())
		fmt.Fprintln(flags.Output(), "Options:")
		flags.PrintDefaults()
	}
	for _, argument := range args {
		if argument == "-h" || argument == "--help" {
			flags.SetOutput(stdout)
			flags.Usage()
			return 0
		}
	}

	orderedArgs, orderErr := reorderInitArgs(args)
	if orderErr != nil {
		fmt.Fprintf(stderr, "omaforge init: %v\n", orderErr)
		return 2
	}
	if err := flags.Parse(orderedArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() > 1 {
		fmt.Fprintln(stderr, "omaforge init accepts exactly one output directory")
		return 2
	}
	if templateName != "bar-widget" {
		fmt.Fprintf(stderr, "unknown template %q; available template: bar-widget\n", templateName)
		return 2
	}

	reader := bufio.NewReader(stdin)
	if flags.NArg() == 1 {
		options.Directory = flags.Arg(0)
	} else if nonInteractive {
		fmt.Fprintln(stderr, "output directory is required in non-interactive mode")
		return 2
	} else {
		options.Directory = prompt(reader, stdout, "Output directory", "")
	}
	if options.Name == "" {
		options.Name = scaffold.DisplayName(options.Directory)
	}
	if !nonInteractive {
		options.Name = prompt(reader, stdout, "Plugin name", options.Name)
		options.ID = prompt(reader, stdout, "Plugin ID", options.ID)
		options.Description = prompt(reader, stdout, "Description", options.Description)
		options.Author = prompt(reader, stdout, "Author", options.Author)
		options.Section = prompt(reader, stdout, "Default section", options.Section)
	}
	if options.Description == "" {
		options.Description = "A polished Omarchy bar widget."
	}
	options.IncludeCI = !noCI
	options.HomeDir, _ = os.UserHomeDir()
	options.OmarchyPath = os.Getenv("OMARCHY_PATH")

	previewed := false
	if options.Force && !options.DryRun {
		previewOptions := options
		previewOptions.DryRun = true
		preview, previewErr := scaffold.Generate(previewOptions)
		if previewErr != nil {
			fmt.Fprintf(stderr, "omaforge init: %v\n", previewErr)
			return 1
		}
		for _, change := range preview.Changes {
			fmt.Fprintf(stdout, "%s %s\n", change.Action, change.Path)
		}
		previewed = true
	}
	result, err := scaffold.Generate(options)
	if err != nil {
		fmt.Fprintf(stderr, "omaforge init: %v\n", err)
		return 1
	}
	if !previewed {
		for _, change := range result.Changes {
			fmt.Fprintf(stdout, "%s %s\n", change.Action, change.Path)
		}
	}
	if options.DryRun {
		fmt.Fprintf(stdout, "Dry run complete: %d files planned; nothing written.\n", len(result.Changes))
	} else {
		fmt.Fprintf(stdout, "Created %s with %d files.\n", result.Directory, len(result.Changes))
		if options.AgentReady {
			fmt.Fprintln(stdout, "Agent-ready files: complete FORGE_SPEC.md, set its status to Ready for implementation, then ask your agent to follow AGENTS.md.")
		}
		fmt.Fprintln(stdout, "Next: omarchy plugin validate "+result.Directory)
	}
	return 0
}

func reorderInitArgs(args []string) ([]string, error) {
	valueFlags := map[string]bool{
		"--id": true, "--name": true, "--description": true, "--author": true,
		"--license": true, "--section": true, "--template": true,
	}
	return reorderCommandArgs(args, valueFlags)
}

func reorderCommandArgs(args []string, valueFlags map[string]bool) ([]string, error) {
	var options, positional []string
	for index := 0; index < len(args); index++ {
		argument := args[index]
		if argument == "--" {
			positional = append(positional, args[index+1:]...)
			break
		}
		if !strings.HasPrefix(argument, "-") || argument == "-" {
			positional = append(positional, argument)
			continue
		}
		options = append(options, argument)
		name := strings.SplitN(argument, "=", 2)[0]
		if valueFlags[name] && !strings.Contains(argument, "=") {
			if index+1 >= len(args) {
				return nil, fmt.Errorf("flag %s requires a value", name)
			}
			index++
			options = append(options, args[index])
		}
	}
	return append(options, positional...), nil
}

func prompt(reader *bufio.Reader, stdout io.Writer, label, defaultValue string) string {
	if defaultValue == "" {
		fmt.Fprintf(stdout, "%s: ", label)
	} else {
		fmt.Fprintf(stdout, "%s [%s]: ", label, defaultValue)
	}
	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultValue
	}
	return value
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
