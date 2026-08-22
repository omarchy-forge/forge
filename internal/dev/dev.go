// Package dev runs a plugin project's explicitly trusted local runtime harness.
package dev

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Run executes the generated isolated runtime harness for directory and state.
func Run(directory, state string, stdin io.Reader, stdout, stderr io.Writer) error {
	return run(directory, state, "", stdin, stdout, stderr)
}

// Screenshot captures a template-declared plugin item without capturing the desktop.
func Screenshot(directory, state, output string, stdin io.Reader, stdout, stderr io.Writer) error {
	if state == "" {
		return fmt.Errorf("screenshot requires a demo state")
	}
	absolute, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("resolve screenshot output: %w", err)
	}
	if filepath.Ext(absolute) != ".png" {
		return fmt.Errorf("screenshot output must end in .png")
	}
	if _, err := os.Lstat(absolute); err == nil {
		return fmt.Errorf("screenshot output already exists: %s", output)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect screenshot output: %w", err)
	}
	if err := run(directory, state, absolute, stdin, stdout, stderr); err != nil {
		return err
	}
	info, err := os.Stat(absolute)
	if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("runtime did not create a nonempty screenshot")
	}
	return nil
}

func run(directory, state, screenshot string, stdin io.Reader, stdout, stderr io.Writer) error {
	if state != "" && state != "ready" && state != "empty" && state != "error" {
		return fmt.Errorf("unsupported demo state %q", state)
	}
	root, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("resolve plugin directory: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return fmt.Errorf("inspect plugin directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("plugin path is not a directory: %s", directory)
	}
	if _, err := os.Stat(filepath.Join(root, "manifest.json")); err != nil {
		return fmt.Errorf("inspect manifest.json: %w", err)
	}
	harness := filepath.Join(root, "tests", "runtime")
	harnessInfo, err := os.Lstat(harness)
	if err != nil {
		return fmt.Errorf("inspect tests/runtime: %w", err)
	}
	if !harnessInfo.Mode().IsRegular() {
		return fmt.Errorf("tests/runtime must be a regular file")
	}

	arguments := []string{harness, "--trust-plugin-code"}
	if state != "" {
		arguments = append(arguments, "--state", state)
	}
	if screenshot != "" {
		arguments = append(arguments, "--screenshot", screenshot)
	}
	command := exec.Command("bash", arguments...)
	command.Dir = root
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("runtime harness: %w", err)
	}
	return nil
}
