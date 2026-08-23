package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeAgentRuntime struct {
	agent     string
	commit    string
	committed bool
	launched  bool
	directory string
}

func (f *fakeAgentRuntime) DefaultAgent(context.Context) (string, error) { return f.agent, nil }
func (f *fakeAgentRuntime) CommitBaseline(directory, _ string) (string, error) {
	f.committed = true
	f.directory = directory
	return f.commit, nil
}
func (f *fakeAgentRuntime) Launch(directory string, _ io.Reader, _, _ io.Writer) error {
	f.launched = true
	f.directory = directory
	return nil
}

func guidedAgentInput(finalAnswer string) string {
	return strings.Join([]string{
		"Agent Widget",
		"my.clock",
		"A local clock with a useful popout.",
		"Example Developer",
		"questionnaire",
		"Current local time.",
		"",
		"Full date and refresh control.",
		"",
		"The local system clock.",
		"",
		"",
		"",
		"",
		"",
		"",
		finalAnswer,
	}, "\n") + "\n"
}

func referenceDrivenAgentInput(finalAnswer string) string {
	return strings.Join([]string{
		"Agent Widget",
		"my.control-center",
		"A local development control center.",
		"Example Developer",
		"references",
		"", // files have already been placed in references/
		"yes",
		"", // no network access
		"", // no persistence
		"", // default failure behavior
		"",
		finalAnswer,
	}, "\n") + "\n"
}

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

func TestRunInitAgentReady(t *testing.T) {
	target := filepath.Join(t.TempDir(), "agent-widget")
	var stdout, stderr bytes.Buffer
	code := Run([]string{
		"init", target,
		"--id", "dev.example.agent-widget",
		"--author", "Example Developer",
		"--description", "A guided agent-built widget.",
		"--agent-ready",
		"--non-interactive",
	}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
	if code != 0 {
		t.Fatalf("Run() code = %d, stderr = %q", code, stderr.String())
	}
	for _, name := range []string{"AGENTS.md", "FORGE_SPEC.md"} {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Fatalf("%s not generated: %v", name, err)
		}
	}
	if !strings.Contains(stdout.String(), "set its status to Ready for implementation") {
		t.Errorf("stdout missing agent-ready guidance: %q", stdout.String())
	}
}

func TestRunInitAgentGuidesCommitsAndLaunchesAfterConfirmation(t *testing.T) {
	target := filepath.Join(t.TempDir(), "agent-widget")
	runtime := &fakeAgentRuntime{agent: "claude", commit: "abc1234"}
	var stdout, stderr bytes.Buffer
	code := runInitWithAgentRuntime(
		[]string{"--agent", target},
		strings.NewReader(guidedAgentInput("yes")),
		&stdout,
		&stderr,
		runtime,
	)
	if code != 0 {
		t.Fatalf("runInitWithAgentRuntime() code = %d, stderr = %q", code, stderr.String())
	}
	if !runtime.committed || !runtime.launched || runtime.directory != target {
		t.Fatalf("runtime = %+v", runtime)
	}
	for _, name := range []string{"AGENTS.md", "FORGE_SPEC.md", "AGENT_PROMPT.md"} {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Fatalf("%s not generated: %v", name, err)
		}
	}
	specification, err := os.ReadFile(filepath.Join(target, "FORGE_SPEC.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(specification, []byte("Specification status: Ready for implementation")) || bytes.Contains(specification, []byte("TODO")) {
		t.Fatalf("guided specification was not completed:\n%s", specification)
	}
	for _, expected := range []string{"Current Omarchy coding agent: claude", "Plugin ID (for example my.clock)", "Baseline commit: abc1234", "Launching claude"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Errorf("stdout missing %q: %q", expected, stdout.String())
		}
	}
}

func TestRunInitAgentCancellationWritesNothing(t *testing.T) {
	target := filepath.Join(t.TempDir(), "agent-widget")
	runtime := &fakeAgentRuntime{agent: "claude", commit: "abc1234"}
	var stdout, stderr bytes.Buffer
	code := runInitWithAgentRuntime(
		[]string{"--agent", target},
		strings.NewReader(guidedAgentInput("cancel")),
		&stdout,
		&stderr,
		runtime,
	)
	if code != 0 {
		t.Fatalf("runInitWithAgentRuntime() code = %d, stderr = %q", code, stderr.String())
	}
	if runtime.committed || runtime.launched {
		t.Fatalf("runtime mutated after cancellation: %+v", runtime)
	}
	entries, err := os.ReadDir(filepath.Join(target, "references"))
	if err != nil || len(entries) != 0 {
		t.Fatalf("cancelled guided init should preserve only the empty reference directory: %v, %v", entries, err)
	}
}

func TestRunInitAgentCopiesConfirmedReferenceBeforeLaunch(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "agent-widget")
	reference := filepath.Join(target, "references", "product brief.md")
	if err := os.MkdirAll(filepath.Dir(reference), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(reference, []byte("A third Omarchy collection project.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runtime := &fakeAgentRuntime{agent: "claude", commit: "abc1234"}
	var stdout, stderr bytes.Buffer
	code := runInitWithAgentRuntime([]string{"--agent", target}, strings.NewReader(referenceDrivenAgentInput("yes")), &stdout, &stderr, runtime)
	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
	content, err := os.ReadFile(filepath.Join(target, "references", "product brief.md"))
	if err != nil || !strings.Contains(string(content), "third Omarchy") {
		t.Fatalf("copied reference = %q, %v", content, err)
	}
	for _, expected := range []string{"Optional references directory created:", "Reference files:", "sha256", "coding agent/provider may receive"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Errorf("stdout missing %q: %s", expected, stdout.String())
		}
	}
	specification, err := os.ReadFile(filepath.Join(target, "FORGE_SPEC.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"extract and implement all visible", "Implement every panel", "Permitted only for functionality required by the confirmed references"} {
		if !strings.Contains(string(specification), expected) {
			t.Errorf("specification missing %q:\n%s", expected, specification)
		}
	}
	prompt, err := os.ReadFile(filepath.Join(target, "AGENT_PROMPT.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(prompt), "implement every mapped item") {
		t.Errorf("prompt does not make reference functionality authoritative:\n%s", prompt)
	}
}

func TestRunInitAgentReferenceModeStopsWhenNoReferenceWasAdded(t *testing.T) {
	target := filepath.Join(t.TempDir(), "agent-widget")
	runtime := &fakeAgentRuntime{agent: "claude", commit: "abc1234"}
	input := strings.Join([]string{
		"Agent Widget",
		"my.control-center",
		"A local development control center.",
		"Example Developer",
		"references",
		"", // finish the empty reference drop step
	}, "\n") + "\n"
	var stdout, stderr bytes.Buffer
	code := runInitWithAgentRuntime([]string{"--agent", target}, strings.NewReader(input), &stdout, &stderr, runtime)
	if code != 0 {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
	if runtime.committed || runtime.launched {
		t.Fatalf("runtime mutated without a reference: %+v", runtime)
	}
	if !strings.Contains(stdout.String(), "Forge found no supported reference files") || !strings.Contains(stdout.String(), "any files already added remain") || !strings.Contains(stdout.String(), "Reference mode requires at least one") || !strings.Contains(stdout.String(), "nothing was written or launched") {
		t.Fatalf("missing safe reference cancellation guidance: %s", stdout.String())
	}
}

func TestRunInitAgentStopsWhenNoAgentIsConfigured(t *testing.T) {
	target := filepath.Join(t.TempDir(), "agent-widget")
	runtime := &fakeAgentRuntime{}
	var stdout, stderr bytes.Buffer
	code := runInitWithAgentRuntime([]string{"--agent", target}, strings.NewReader(""), &stdout, &stderr, runtime)
	if code == 0 || !strings.Contains(stderr.String(), "omarchy default agent <name>") {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("missing-agent preflight wrote target: %v", err)
	}
}

func TestRunInitAgentRejectsUnsafeModeCombinations(t *testing.T) {
	for _, argument := range []string{"--non-interactive", "--dry-run", "--force", "--agent-ready"} {
		var stdout, stderr bytes.Buffer
		code := runInitWithAgentRuntime(
			[]string{"--agent", argument, "agent-widget"},
			strings.NewReader(""),
			&stdout,
			&stderr,
			&fakeAgentRuntime{agent: "claude"},
		)
		if code != 2 || !strings.Contains(stderr.String(), "cannot be combined") {
			t.Errorf("argument %s: code = %d, stderr = %q", argument, code, stderr.String())
		}
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
	if !strings.Contains(stdout.String(), "-agent-ready") {
		t.Errorf("stdout missing agent-ready option: %q", stdout.String())
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
	if !strings.Contains(stdout.String(), "Plugin ID (for example my.clock):") || !strings.Contains(stdout.String(), "nothing written") {
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
