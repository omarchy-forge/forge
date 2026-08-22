package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDevRequiresTrustAcknowledgement(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"dev", "."}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 2 || !strings.Contains(stderr.String(), "--trust-plugin-code") {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
}

func TestDevRunsProjectHarness(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte("echo dev-ok\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{"dev", root, "--trust-plugin-code"}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 || stdout.String() != "dev-ok\n" || stderr.Len() != 0 {
		t.Fatalf("code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
	}
}

func TestDevAcceptsTrustFlagBeforeDirectory(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte("exit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{"dev", "--trust-plugin-code", root}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
}

func TestDevPassesValidatedState(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "tests"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("{}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tests", "runtime"), []byte("printf '%s\\n' \"$@\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{"dev", root, "--state", "error", "--trust-plugin-code"}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 || stdout.String() != "--trust-plugin-code\n--state\nerror\n" {
		t.Fatalf("code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
	}
}

func TestDevRejectsUnknownState(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"dev", ".", "--trust-plugin-code", "--state", "loading"}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 2 || !strings.Contains(stderr.String(), "unknown state") {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
}
