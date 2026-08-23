package checks

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestRunValidProjectIsDeterministic(t *testing.T) {
	root := validProject(t)
	first := Run(root)
	second := Run(root)
	if !reflect.DeepEqual(first, second) {
		t.Fatal("reports differ for identical input")
	}
	if first.Summary.Errors != 0 {
		t.Fatalf("errors = %d, findings = %#v", first.Summary.Errors, first.Findings)
	}
	if first.Compatibility != "omarchy-4" {
		t.Errorf("compatibility = %q", first.Compatibility)
	}
}

func TestRunReportsOfficialAndForgeRules(t *testing.T) {
	root := t.TempDir()
	write(t, root, "manifest.json", `{"schemaVersion":"1","id":"omarchy.bad","name":"Bad","version":"latest","kinds":["mystery"],"entryPoints":{}}`)
	report := Run(root)
	ids := map[string]bool{}
	sources := map[string]bool{}
	for _, f := range report.Findings {
		ids[f.RuleID] = true
		sources[f.Source] = true
	}
	for _, id := range []string{"OF102", "OF104", "OF107", "OF200", "OF201"} {
		if !ids[id] {
			t.Errorf("missing %s in %#v", id, report.Findings)
		}
	}
	if !sources["official-parity"] || !sources["forge"] {
		t.Errorf("sources = %v", sources)
	}
}

func TestRenderersAreMachineReadable(t *testing.T) {
	report := Run(validProject(t))
	var output bytes.Buffer
	if err := WriteJSON(&output, report); err != nil {
		t.Fatal(err)
	}
	var decoded Report
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	output.Reset()
	if err := WriteSARIF(&output, report); err != nil {
		t.Fatal(err)
	}
	var sarif map[string]any
	if err := json.Unmarshal(output.Bytes(), &sarif); err != nil {
		t.Fatal(err)
	}
	if sarif["version"] != "2.1.0" {
		t.Errorf("SARIF version = %v", sarif["version"])
	}
	output.Reset()
	if err := WriteText(&output, report); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "schema 1, omarchy-4") {
		t.Errorf("text = %q", output.String())
	}
}

func TestRunReportsSymlinkAndSourceHeuristics(t *testing.T) {
	root := validProject(t)
	if err := os.Symlink("Panel.qml", filepath.Join(root, "linked.qml")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	write(t, root, "unsafe.qml", "import QtQuick\nItem { property color bad: \"#ff00ff\"; property var command: \"sudo pacman -S thing\" }\n")
	report := Run(root)
	ids := map[string]bool{}
	for _, f := range report.Findings {
		ids[f.RuleID] = true
	}
	for _, id := range []string{"OF106", "OF301", "OF302", "OF303", "OF304"} {
		if !ids[id] {
			t.Errorf("missing %s in %#v", id, report.Findings)
		}
	}
}

func TestRunRejectsNULByteInSource(t *testing.T) {
	root := validProject(t)
	path := filepath.Join(root, "services", "DataService.qml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("import QtQuick\nItem { property string separator: \"\x00\" }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	report := Run(root)
	for _, finding := range report.Findings {
		if finding.RuleID == "OF305" {
			if finding.Severity != Error || finding.Path != "services/DataService.qml" || finding.Line != 2 {
				t.Fatalf("unexpected OF305 finding: %#v", finding)
			}
			return
		}
	}
	t.Fatalf("missing OF305 in %#v", report.Findings)
}

func validProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	write(t, root, "manifest.json", `{"schemaVersion":1,"id":"dev.example.good","name":"Good","version":"0.1.0","kinds":["bar-widget"],"entryPoints":{"barWidget":"Panel.qml"}}`)
	write(t, root, "Panel.qml", "import QtQuick\nItem {}\n")
	write(t, root, "LICENSE", "MIT\n")
	write(t, root, "assets/preview.png", "not-an-image")
	write(t, root, "README.md", "# Good\n\nCompatible with Omarchy 4.\n\n## Requirements\n## Install\n## Configuration\n## Development\n## Update\n## Removal\n## Privacy\n## License\n")
	return root
}

func write(t *testing.T, root, name, contents string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
