package scaffold

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrepareReferencesValidatesAndNamesDeterministically(t *testing.T) {
	directory := t.TempDir()
	textPath := filepath.Join(directory, "product brief.md")
	imagePath := filepath.Join(directory, "mockup.PNG")
	textContent := []byte("Build a quiet local clock.\n")
	imageContent := []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 1}
	if err := os.WriteFile(textPath, textContent, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(imagePath, imageContent, 0o644); err != nil {
		t.Fatal(err)
	}

	references, err := PrepareReferences([]string{textPath, imagePath})
	if err != nil {
		t.Fatal(err)
	}
	if references[0].ProjectPath != "references/reference-01.md" || references[0].Kind != "text" {
		t.Fatalf("text reference = %+v", references[0])
	}
	if references[1].ProjectPath != "references/reference-02.png" || references[1].Kind != "image" {
		t.Fatalf("image reference = %+v", references[1])
	}
	digest := sha256.Sum256(textContent)
	if references[0].SHA256 != hex.EncodeToString(digest[:]) {
		t.Fatalf("digest = %s", references[0].SHA256)
	}
}

func TestPrepareReferencesAcceptsValidatedSVGWithoutRenderingIt(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "app-icon.svg")
	content := []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><path d="M0 0h16v16z"/></svg>`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	references, err := PrepareReferences([]string{path})
	if err != nil {
		t.Fatal(err)
	}
	if references[0].Kind != "image" || references[0].ProjectPath != "references/reference-01.svg" {
		t.Fatalf("SVG reference = %+v", references[0])
	}

	if err := os.WriteFile(path, []byte(`<html/>`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareReferences([]string{path}); err == nil || !strings.Contains(err.Error(), "svg root") {
		t.Fatalf("invalid SVG error = %v", err)
	}
}

func TestPrepareReferencesRejectsUnsafeOrInvalidInputs(t *testing.T) {
	directory := t.TempDir()
	valid := filepath.Join(directory, "valid.txt")
	if err := os.WriteFile(valid, []byte("brief"), 0o644); err != nil {
		t.Fatal(err)
	}
	symlink := filepath.Join(directory, "link.txt")
	if err := os.Symlink(valid, symlink); err != nil {
		t.Fatal(err)
	}
	cases := map[string]struct {
		path     string
		contains string
	}{
		"symlink":   {symlink, "regular file"},
		"extension": {filepath.Join(directory, "brief.pdf"), "supported extensions"},
		"duplicate": {valid, "selected more than once"},
	}
	if err := os.WriteFile(cases["extension"].path, []byte("pdf"), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, test := range cases {
		paths := []string{test.path}
		if name == "duplicate" {
			paths = []string{valid, valid}
		}
		_, err := PrepareReferences(paths)
		if err == nil || !strings.Contains(err.Error(), test.contains) {
			t.Errorf("%s error = %v", name, err)
		}
	}
	badImage := filepath.Join(directory, "bad.png")
	if err := os.WriteFile(badImage, []byte("not png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareReferences([]string{badImage}); err == nil || !strings.Contains(err.Error(), "signature") {
		t.Errorf("bad image error = %v", err)
	}
	badText := filepath.Join(directory, "bad.md")
	if err := os.WriteFile(badText, []byte{0xff}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareReferences([]string{badText}); err == nil || !strings.Contains(err.Error(), "UTF-8") {
		t.Errorf("bad text error = %v", err)
	}
}

func TestGenerateCopiesPreparedReferences(t *testing.T) {
	directory := t.TempDir()
	source := filepath.Join(directory, "brief.txt")
	content := []byte("third collection project\n")
	if err := os.WriteFile(source, content, 0o644); err != nil {
		t.Fatal(err)
	}
	references, err := PrepareReferences([]string{source})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(directory, "agent-widget")
	options := validOptions(target)
	options.AgentReady, options.AgentGuided, options.ReferenceDriven = true, true, true
	options.BarSummary, options.ClickBehavior = "Status", "Open details"
	options.PopoutSummary, options.DataSources = "Details", "Local data"
	options.UserActions = "Refresh details"
	options.LocalCommands, options.NetworkAccess = "Not required", "Not required"
	options.Persistence, options.FailureBehavior = "Not required", "Show an error"
	options.References = references
	if _, err := Generate(options); err != nil {
		t.Fatal(err)
	}
	written, err := os.ReadFile(filepath.Join(target, "references", "reference-01.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(written) != string(content) {
		t.Fatalf("reference content = %q", written)
	}
	prompt, err := os.ReadFile(filepath.Join(target, "AGENT_PROMPT.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		references[0].SHA256,
		"untrusted data",
		"requirement inventory",
		"every visible function",
		"Map every inventory item",
		"decoration or styling inspiration only",
		"stop and ask the user to resolve the specification",
		"Supplied logos and icons are required",
		"reference-versus-render visual review",
	} {
		if !strings.Contains(string(prompt), required) {
			t.Fatalf("prompt missing reference requirement %q:\n%s", required, prompt)
		}
	}
	specification, err := os.ReadFile(filepath.Join(target, "FORGE_SPEC.md"))
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"Visual fidelity contract", "Raster assets do not need to be recreated as vectors", "visually match the confirmed design references"} {
		if !strings.Contains(string(specification), required) {
			t.Fatalf("specification missing visual-fidelity requirement %q:\n%s", required, specification)
		}
	}
}
