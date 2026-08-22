package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCheckFormatsAndExitCodes(t *testing.T) {
	root := t.TempDir()
	manifest := `{"schemaVersion":1,"id":"dev.example.test","name":"Test","version":"0.1.0","kinds":["bar-widget"],"entryPoints":{"barWidget":"Panel.qml"}}`
	for name, contents := range map[string]string{"manifest.json": manifest, "Panel.qml": "import QtQuick\nItem {}", "LICENSE": "MIT", "assets/preview.png": "x", "README.md": "# Test\nCompatible with Omarchy 4\n## Requirements\n## Install\n## Configuration\n## Development\n## Update\n## Removal\n## Privacy\n## License\n"} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, format := range []string{"text", "json", "sarif"} {
		t.Run(format, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := Run([]string{"check", root, "--format", format}, strings.NewReader(""), &stdout, &stderr, BuildInfo{})
			if code != 0 {
				t.Fatalf("code=%d stderr=%q output=%q", code, stderr.String(), stdout.String())
			}
			if stdout.Len() == 0 || stderr.Len() != 0 {
				t.Errorf("stdout=%q stderr=%q", stdout.String(), stderr.String())
			}
		})
	}
	if code := Run([]string{"check", t.TempDir()}, strings.NewReader(""), ioDiscard{}, ioDiscard{}, BuildInfo{}); code != 1 {
		t.Errorf("invalid project code=%d, want 1", code)
	}
	if code := Run([]string{"check", root, "--format", "xml"}, strings.NewReader(""), ioDiscard{}, ioDiscard{}, BuildInfo{}); code != 2 {
		t.Errorf("bad format code=%d, want 2", code)
	}
	if code := Run([]string{"check", root, "--omarchy-version", "5"}, strings.NewReader(""), ioDiscard{}, ioDiscard{}, BuildInfo{}); code != 2 {
		t.Errorf("unsupported compatibility code=%d, want 2", code)
	}
}
