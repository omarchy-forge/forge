package scaffold

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGenerateCreatesDeterministicProject(t *testing.T) {
	first := filepath.Join(t.TempDir(), "project-pulse")
	second := filepath.Join(t.TempDir(), "project-pulse")
	options := validOptions(first)
	if _, err := Generate(options); err != nil {
		t.Fatalf("Generate(first): %v", err)
	}
	options.Directory = second
	if _, err := Generate(options); err != nil {
		t.Fatalf("Generate(second): %v", err)
	}

	firstSnapshot := snapshot(t, first)
	secondSnapshot := snapshot(t, second)
	if !reflect.DeepEqual(firstSnapshot, secondSnapshot) {
		t.Errorf("generated trees differ\nfirst:  %v\nsecond: %v", firstSnapshot, secondSnapshot)
	}

	wantGolden := readGolden(t, "testdata/bar-widget.golden")
	if !reflect.DeepEqual(firstSnapshot, wantGolden) {
		t.Errorf("generated tree differs from golden\ngot:  %v\nwant: %v", firstSnapshot, wantGolden)
	}

	manifestBytes, err := os.ReadFile(filepath.Join(first, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest map[string]any
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		t.Fatalf("manifest is invalid JSON: %v", err)
	}
	if manifest["id"] != "dev.example.project-pulse" {
		t.Errorf("manifest id = %v", manifest["id"])
	}
	for _, executable := range []string{"demo/run", "tests/run", "tests/runtime"} {
		info, err := os.Stat(filepath.Join(first, executable))
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o755 {
			t.Errorf("%s mode = %o, want 755", executable, info.Mode().Perm())
		}
	}
	demoScript, err := os.ReadFile(filepath.Join(first, "demo", "run"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"omarchy-shell", "setDemoState", "Applied %s without changing shell configuration"} {
		if !bytes.Contains(demoScript, []byte(required)) {
			t.Errorf("demo/run missing %q", required)
		}
	}
	for _, forbidden := range []string{"OMAFORGE_DEMO_STATE", "omarchy-restart-shell"} {
		if bytes.Contains(demoScript, []byte(forbidden)) {
			t.Errorf("demo/run contains obsolete runtime mechanism %q", forbidden)
		}
	}
	runtimeScript, err := os.ReadFile(filepath.Join(first, "tests", "runtime"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"--trust-plugin-code", "--state", "XDG_RUNTIME_DIR", "OMAFORGE_DEV_STATE", "OMAFORGE_RUNTIME_PASS"} {
		if !bytes.Contains(runtimeScript, []byte(required)) {
			t.Errorf("tests/runtime missing %q", required)
		}
	}
	for _, forbidden := range []string{"omarchy plugin add", "omarchy plugin enable", "omarchy-shell"} {
		if bytes.Contains(runtimeScript, []byte(forbidden)) {
			t.Errorf("tests/runtime contains persistent-shell operation %q", forbidden)
		}
	}
}

func TestGenerateDryRunWritesNothing(t *testing.T) {
	target := filepath.Join(t.TempDir(), "dry-run")
	options := validOptions(target)
	options.DryRun = true
	result, err := Generate(options)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changes) == 0 {
		t.Fatal("dry run returned no planned changes")
	}
	if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("dry run created target: %v", err)
	}
}

func TestGenerateCanOmitCIAndInitializeGit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	target := filepath.Join(t.TempDir(), "local-widget")
	options := validOptions(target)
	options.IncludeCI = false
	options.InitGit = true
	if _, err := Generate(options); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(target, ".github", "workflows", "forge.yml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("CI workflow exists with IncludeCI=false: %v", err)
	}
	command := exec.Command("git", "-C", target, "branch", "--show-current")
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(output)) != "main" {
		t.Errorf("initial branch = %q, want main", output)
	}
}

func TestGenerateRefusesUnsafeTargets(t *testing.T) {
	temporary := t.TempDir()
	fileTarget := filepath.Join(temporary, "plain-file")
	if err := os.WriteFile(fileTarget, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	symlinkTarget := filepath.Join(temporary, "linked-project")
	if err := os.Symlink(temporary, symlinkTarget); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		target  string
		mutate  func(*Options)
		message string
	}{
		{name: "root", target: string(filepath.Separator), message: "filesystem root"},
		{name: "home", target: temporary, mutate: func(o *Options) { o.HomeDir = temporary }, message: "home directory"},
		{name: "omarchy", target: temporary, mutate: func(o *Options) { o.OmarchyPath = temporary }, message: "Omarchy installation"},
		{name: "inside omarchy", target: filepath.Join(temporary, "owned", "plugin"), mutate: func(o *Options) { o.OmarchyPath = filepath.Join(temporary, "owned") }, message: "Omarchy installation"},
		{name: "file", target: fileTarget, message: "not a directory"},
		{name: "symlink", target: symlinkTarget, message: "symlink"},
		{name: "bad slug", target: filepath.Join(temporary, "Bad Name"), message: "directory name"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := validOptions(test.target)
			if test.mutate != nil {
				test.mutate(&options)
			}
			_, err := Generate(options)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("Generate() error = %v, want substring %q", err, test.message)
			}
		})
	}
}

func TestGenerateRefusesSymlinkAncestor(t *testing.T) {
	root := t.TempDir()
	linked := filepath.Join(root, "linked")
	if err := os.Symlink(t.TempDir(), linked); err != nil {
		t.Fatal(err)
	}
	options := validOptions(filepath.Join(linked, "project-pulse"))
	_, err := Generate(options)
	if err == nil || !strings.Contains(err.Error(), "contains a symlink") {
		t.Fatalf("Generate() error = %v, want symlink ancestor refusal", err)
	}
}

func TestGenerateRefusesInvalidValues(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Options)
		message string
	}{
		{name: "unqualified id", mutate: func(o *Options) { o.ID = "clock" }, message: "namespaced"},
		{name: "reserved id", mutate: func(o *Options) { o.ID = "omarchy.clock" }, message: "omarchy.*"},
		{name: "traversal id", mutate: func(o *Options) { o.ID = "dev..clock" }, message: "namespaced"},
		{name: "newline name", mutate: func(o *Options) { o.Name = "Bad\nName" }, message: "single line"},
		{name: "unsupported license", mutate: func(o *Options) { o.License = "GPL-3.0" }, message: "unsupported license"},
		{name: "bad section", mutate: func(o *Options) { o.Section = "bottom" }, message: "must be left"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := validOptions(filepath.Join(t.TempDir(), "project-pulse"))
			test.mutate(&options)
			_, err := Generate(options)
			if err == nil || !strings.Contains(err.Error(), test.message) {
				t.Fatalf("Generate() error = %v, want substring %q", err, test.message)
			}
		})
	}
}

func TestGenerateCollisionRequiresForceAndPreservesUnrelatedFiles(t *testing.T) {
	target := filepath.Join(t.TempDir(), "project-pulse")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	unrelated := filepath.Join(target, "notes.txt")
	if err := os.WriteFile(unrelated, []byte("user content"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifest := filepath.Join(target, "manifest.json")
	if err := os.WriteFile(manifest, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	options := validOptions(target)
	if _, err := Generate(options); err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("Generate() error = %v, want nonempty refusal", err)
	}
	options.Force = true
	result, err := Generate(options)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changes[changeIndex(t, result.Changes, manifest)].Action != "OVERWRITE" {
		t.Errorf("manifest action was not OVERWRITE: %v", result.Changes)
	}
	content, err := os.ReadFile(unrelated)
	if err != nil || string(content) != "user content" {
		t.Fatalf("unrelated file changed: content=%q err=%v", content, err)
	}
}

func TestGenerateAcceptsGitOnlyTargetWithoutForce(t *testing.T) {
	target := filepath.Join(t.TempDir(), "project-pulse")
	gitDirectory := filepath.Join(target, ".git")
	if err := os.MkdirAll(gitDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(gitDirectory, "HEAD")
	if err := os.WriteFile(marker, []byte("ref: refs/heads/main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Generate(validOptions(target)); err != nil {
		t.Fatalf("Generate() error = %v, want .git-only target to be accepted", err)
	}
	content, err := os.ReadFile(marker)
	if err != nil || string(content) != "ref: refs/heads/main\n" {
		t.Fatalf("Git metadata changed: content=%q err=%v", content, err)
	}
}

func TestGenerateRefusesGitMetadataFileWithoutForce(t *testing.T) {
	target := filepath.Join(t.TempDir(), "project-pulse")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, ".git"), []byte("gitdir: elsewhere\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Generate(validOptions(target))
	if err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("Generate() error = %v, want nonempty refusal", err)
	}
}

func TestGenerateForceRefusesNestedSymlink(t *testing.T) {
	target := filepath.Join(t.TempDir(), "project-pulse")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(target, "components")); err != nil {
		t.Fatal(err)
	}
	options := validOptions(target)
	options.Force = true
	_, err := Generate(options)
	if err == nil || !strings.Contains(err.Error(), "contains a symlink") {
		t.Fatalf("Generate() error = %v, want nested symlink refusal", err)
	}
}

func TestGenerateEscapesUserContent(t *testing.T) {
	target := filepath.Join(t.TempDir(), "safe-project")
	options := validOptions(target)
	options.Name = `Name *with* [markup] "quotes"`
	options.Description = `Description with <tags> & symbols.`
	options.Author = `Author *Name*`
	if _, err := Generate(options); err != nil {
		t.Fatal(err)
	}
	manifest, err := os.ReadFile(filepath.Join(target, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(manifest, &decoded); err != nil {
		t.Fatalf("manifest JSON: %v", err)
	}
	if decoded["name"] != options.Name {
		t.Errorf("manifest name = %q, want %q", decoded["name"], options.Name)
	}
	readme, err := os.ReadFile(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(readme, []byte(`Name \*with\* \[markup\]`)) {
		t.Errorf("README did not escape Markdown: %s", readme)
	}
}

func TestGeneratedPluginPassesOfficialValidator(t *testing.T) {
	if _, err := exec.LookPath("omarchy"); err != nil {
		t.Skip("omarchy is not installed")
	}
	target := filepath.Join(t.TempDir(), "project-pulse")
	if _, err := Generate(validOptions(target)); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("omarchy", "plugin", "validate", target)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("official validator failed: %v\n%s", err, output)
	}
}

func validOptions(directory string) Options {
	return Options{
		Directory:   directory,
		ID:          "dev.example.project-pulse",
		Name:        "Project Pulse",
		Description: "Active project status at a glance.",
		Author:      "Example Developer",
		License:     "MIT",
		Section:     "right",
		IncludeCI:   true,
		HomeDir:     filepath.Join(filepath.Dir(directory), "not-home"),
		OmarchyPath: "/usr/share/omarchy",
	}
}

func snapshot(t *testing.T, root string) []string {
	t.Helper()
	var entries []string
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		content, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		hash := sha256.Sum256(content)
		relative, _ := filepath.Rel(root, name)
		entries = append(entries, filepath.ToSlash(relative)+" "+info.Mode().Perm().String()+" "+hex.EncodeToString(hash[:]))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(entries)
	return entries
}

func readGolden(t *testing.T, name string) []string {
	t.Helper()
	content, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	sort.Strings(lines)
	return lines
}

func changeIndex(t *testing.T, changes []Change, path string) int {
	t.Helper()
	for index, change := range changes {
		if change.Path == path {
			return index
		}
	}
	t.Fatalf("no change for %s", path)
	return -1
}
