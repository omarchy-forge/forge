package dev

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunExecutesTrustedHarness(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	script := "printf '%s\\n' \"$1\"\nprintf '%s\\n' \"$PWD\"\n"
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := Run(root, "ready", strings.NewReader(""), &output, &output); err != nil {
		t.Fatal(err)
	}
	want := "--trust-plugin-code\n" + root + "\n"
	if output.String() != want {
		t.Fatalf("output = %q, want %q", output.String(), want)
	}
}

func TestWatchRerunsAfterChangeAndStops(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	counter := filepath.Join(t.TempDir(), "runs")
	script := "printf x >> " + counter + "\n"
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- Watch(ctx, root, "ready", strings.NewReader(""), io.Discard, io.Discard) }()
	waitForSize(t, counter, 1)
	if err := os.WriteFile(filepath.Join(root, "Panel.qml"), []byte("Item {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	waitForSize(t, counter, 2)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watch did not stop")
	}
}

func waitForSize(t *testing.T, path string, minimum int64) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if info, err := os.Stat(path); err == nil && info.Size() >= minimum {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("%s did not reach size %d", path, minimum)
}

func TestRunRejectsMissingHarness(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Run(root, "", strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "tests/runtime") {
		t.Fatalf("error = %v, want missing runtime harness", err)
	}
}

func TestRunRejectsSymlinkHarness(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "runtime")
	if err := os.WriteFile(outside, []byte("exit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "tests", "runtime")); err != nil {
		t.Fatal(err)
	}
	err := Run(root, "", strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("error = %v, want regular-file rejection", err)
	}
}

func TestRunRejectsUnknownState(t *testing.T) {
	err := Run(".", "loading", strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "unsupported demo state") {
		t.Fatalf("error = %v, want unsupported state", err)
	}
}

func TestScreenshotPassesAbsoluteOutput(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte("printf '%s\\n' \"$@\"\nprintf image > \"${@: -1}\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "preview.png")
	var buffer bytes.Buffer
	if err := Screenshot(root, "ready", output, strings.NewReader(""), &buffer, &buffer); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buffer.String(), "--screenshot\n"+output) {
		t.Fatalf("arguments = %q", buffer.String())
	}
}

func TestScreenshotRefusesOverwrite(t *testing.T) {
	output := filepath.Join(t.TempDir(), "preview.png")
	if err := os.WriteFile(output, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Screenshot(".", "ready", output, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("error = %v", err)
	}
}
