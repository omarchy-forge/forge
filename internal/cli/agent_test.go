package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecAgentRuntimeCommitsGeneratedBaseline(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is unavailable")
	}
	directory := filepath.Join(t.TempDir(), "agent-widget")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	for key, value := range map[string]string{
		"GIT_AUTHOR_NAME": "Forge Test", "GIT_AUTHOR_EMAIL": "forge@example.invalid",
		"GIT_COMMITTER_NAME": "Forge Test", "GIT_COMMITTER_EMAIL": "forge@example.invalid",
	} {
		t.Setenv(key, value)
	}
	if output, err := exec.Command("git", "init", "--initial-branch=main", directory).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, output)
	}
	if err := os.WriteFile(filepath.Join(directory, "README.md"), []byte("generated\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	commit, err := (execAgentRuntime{}).CommitBaseline(directory, "Agent Widget")
	if err != nil {
		t.Fatal(err)
	}
	if len(commit) < 7 {
		t.Fatalf("short commit = %q", commit)
	}
	output, err := exec.Command("git", "-C", directory, "log", "-1", "--format=%s").Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(output)) != "chore: scaffold Agent Widget with Omarchy Forge" {
		t.Fatalf("commit subject = %q", output)
	}
}
