package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitNonInteractive(t *testing.T) {
	target := filepath.Join(t.TempDir(), "my-widget")
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"init", target,
		"--id", "dev.example.my-widget",
		"--author", "Example Developer",
		"--non-interactive",
	}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(target, "manifest.json")); err != nil {
		t.Fatalf("manifest not generated: %v", err)
	}
	if !strings.Contains(stdout.String(), "Created ") || !strings.Contains(stdout.String(), "omarchy plugin validate") {
		t.Errorf("stdout missing completion guidance: %q", stdout.String())
	}
}

func TestRunInitAcceptsGitOnlyTarget(t *testing.T) {
	target := filepath.Join(t.TempDir(), "my-widget")
	if err := os.MkdirAll(filepath.Join(target, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"init", target,
		"--id", "dev.example.my-widget",
		"--author", "Example Developer",
		"--non-interactive",
	}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(target, "manifest.json")); err != nil {
		t.Fatalf("manifest not generated: %v", err)
	}
}

func TestRunInitHelpUsesStdout(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", "--help"}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage: omaforge init") {
		t.Errorf("stdout missing init usage: %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunInitDryRun(t *testing.T) {
	target := filepath.Join(t.TempDir(), "my-widget")
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"init", target,
		"--id", "dev.example.my-widget",
		"--author", "Example Developer",
		"--non-interactive",
		"--dry-run",
	}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote target: %v", err)
	}
	if !strings.Contains(stdout.String(), "nothing written") {
		t.Errorf("stdout missing dry-run result: %q", stdout.String())
	}
}

func TestRunInitInteractivePrompts(t *testing.T) {
	target := filepath.Join(t.TempDir(), "prompted-widget")
	input := strings.Join([]string{
		target,
		"Prompted Widget",
		"dev.example.prompted-widget",
		"A prompted widget.",
		"Example Developer",
		"left",
	}, "\n") + "\n"
	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", "--dry-run"}, strings.NewReader(input), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Plugin ID:") || !strings.Contains(stdout.String(), "nothing written") {
		t.Errorf("stdout missing prompts or result: %q", stdout.String())
	}
}

func TestRunInitMissingNonInteractiveValues(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Run([]string{"init", "my-widget", "--non-interactive"}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code == 0 {
		t.Fatal("Run() code = 0, want nonzero")
	}
	if !strings.Contains(stderr.String(), "plugin ID") {
		t.Errorf("stderr = %q, want plugin ID error", stderr.String())
	}
}

func TestRunInitForcePrintsOverwriteBeforeCompletion(t *testing.T) {
	target := filepath.Join(t.TempDir(), "my-widget")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "manifest.json"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"init", target, "--id", "dev.example.my-widget", "--author", "Example Developer",
		"--non-interactive", "--force",
	}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	preview := strings.Index(stdout.String(), "OVERWRITE "+filepath.Join(target, "manifest.json"))
	completion := strings.Index(stdout.String(), "Created ")
	if preview < 0 || completion < 0 || preview > completion {
		t.Errorf("overwrite preview did not precede completion: %q", stdout.String())
	}
}
