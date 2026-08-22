// Package dev runs a plugin project's explicitly trusted local runtime harness.
package dev

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

// Run executes the generated isolated runtime harness for directory and state.
func Run(directory, state string, stdin io.Reader, stdout, stderr io.Writer) error {
	return run(context.Background(), directory, state, "", stdin, stdout, stderr)
}

// Watch reruns the isolated harness after debounced project changes until ctx ends.
func Watch(ctx context.Context, directory, state string, stdin io.Reader, stdout, stderr io.Writer) error {
	root, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("resolve plugin directory: %w", err)
	}
	if err := run(ctx, root, state, "", stdin, stdout, stderr); err != nil && ctx.Err() == nil {
		fmt.Fprintf(stderr, "omaforge dev: %v; watching for changes\n", err)
	}
	if ctx.Err() != nil {
		fmt.Fprintln(stdout, "Stopped watching.")
		return nil
	}
	fingerprint, err := projectFingerprint(root)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, "Watching for plugin changes. Press Ctrl-C to stop.")
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(stdout, "Stopped watching.")
			return nil
		case <-ticker.C:
			next, fingerprintErr := projectFingerprint(root)
			if fingerprintErr != nil {
				fmt.Fprintf(stderr, "omaforge dev: scan changes: %v\n", fingerprintErr)
				continue
			}
			if next == fingerprint {
				continue
			}
			fingerprint = next
			fmt.Fprintln(stdout, "Plugin change detected; rerunning isolated harness.")
			if runErr := run(ctx, root, state, "", stdin, stdout, stderr); runErr != nil && ctx.Err() == nil {
				fmt.Fprintf(stderr, "omaforge dev: %v; continuing to watch\n", runErr)
			}
		}
	}
}

func projectFingerprint(root string) ([32]byte, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && entry.Name() == ".git" {
			return filepath.SkipDir
		}
		if entry.Type().IsRegular() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return [32]byte{}, fmt.Errorf("scan plugin directory: %w", err)
	}
	sort.Strings(paths)
	hash := sha256.New()
	for _, path := range paths {
		relative, _ := filepath.Rel(root, path)
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			return [32]byte{}, fmt.Errorf("read %s: %w", relative, readErr)
		}
		fmt.Fprintf(hash, "%s\x00", filepath.ToSlash(relative))
		hash.Write(contents)
	}
	var result [32]byte
	copy(result[:], hash.Sum(nil))
	return result, nil
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
	if err := run(context.Background(), directory, state, absolute, stdin, stdout, stderr); err != nil {
		return err
	}
	info, err := os.Stat(absolute)
	if err != nil || !info.Mode().IsRegular() || info.Size() == 0 {
		return fmt.Errorf("runtime did not create a nonempty screenshot")
	}
	return nil
}

func run(ctx context.Context, directory, state, screenshot string, stdin io.Reader, stdout, stderr io.Writer) error {
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
	command := exec.CommandContext(ctx, "bash", arguments...)
	command.Dir = root
	command.Stdin = stdin
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("runtime harness: %w", err)
	}
	return nil
}
