package doctor

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeRunner struct {
	missing  map[string]bool
	outputs  map[string]string
	failures map[string]bool
}

func (f fakeRunner) LookPath(name string) (string, error) {
	if f.missing[name] {
		return "", errors.New("missing")
	}
	return "/bin/" + name, nil
}
func (f fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	if f.failures[key] {
		return []byte(f.outputs[key]), errors.New("failed")
	}
	return []byte(f.outputs[key]), nil
}

func TestRunCombinesEnvironmentAndProjectChecks(t *testing.T) {
	root := t.TempDir()
	write := func(name, contents string) {
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("manifest.json", `{"schemaVersion":1,"id":"dev.example.test","name":"Test","version":"0.1.0","kinds":["bar-widget"],"entryPoints":{"barWidget":"Panel.qml"}}`)
	write("Panel.qml", "import QtQuick\nItem{}")
	runner := fakeRunner{missing: map[string]bool{"qmllint": true}, outputs: map[string]string{"pacman -Q omarchy": "omarchy 4.0.0-1\n", "omarchy-shell shell ping": "ok\n"}, failures: map[string]bool{}}
	report := Run(root, runner)
	if !report.Failed() {
		t.Fatal("missing README/LICENSE project should fail")
	}
	statuses := map[string]Status{}
	for _, d := range report.Diagnostics {
		statuses[d.RuleID] = d.Status
	}
	if statuses["OD100"] != Pass || statuses["OD103"] != Pass || statuses["OD104"] != Warn || statuses["OD105"] != Pass {
		t.Errorf("statuses=%v", statuses)
	}
}

func TestRunReportsUnavailableEnvironment(t *testing.T) {
	runner := fakeRunner{missing: map[string]bool{"omarchy": true, "quickshell": true, "omarchy-shell": true, "qmllint": true}, outputs: map[string]string{}, failures: map[string]bool{}}
	report := Run(t.TempDir(), runner)
	if !report.Failed() {
		t.Fatal("unavailable environment should fail")
	}
}
