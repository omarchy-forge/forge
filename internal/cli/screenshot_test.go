package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScreenshotRunsHarness(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte("printf image > \"${@: -1}\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "preview.png")
	var stdout, stderr bytes.Buffer
	code := Run([]string{"screenshot", root, "--trust-plugin-code", "--state", "ready", "--output", output}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 || !strings.Contains(stdout.String(), "Saved plugin-only screenshot") {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestScreenshotRequiresStateAndTrust(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"screenshot", ".", "--output", "preview.png"}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 2 {
		t.Fatalf("code=%d stderr=%q", code, stderr.String())
	}
}
