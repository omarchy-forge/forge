// Package dev runs a plugin project's explicitly trusted local runtime harness.
package dev

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Run executes the generated isolated runtime harness for directory.
func Run(directory string, stdin io.Reader, stdout, stderr io.Writer) error {
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

	command := exec.Command("bash", harness, "--trust-plugin-code")
	command.Dir = root
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("runtime harness: %w", err)
	}
	return nil
}
