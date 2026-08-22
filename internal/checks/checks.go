// Package checks implements deterministic, non-executing plugin validation.
package checks

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const SchemaVersion = "1"

type Severity string

const (
	Error   Severity = "error"
	Warning Severity = "warning"
	Note    Severity = "note"
)

type Finding struct {
	RuleID      string   `json:"ruleId"`
	Source      string   `json:"source"`
	Severity    Severity `json:"severity"`
	Message     string   `json:"message"`
	Path        string   `json:"path,omitempty"`
	Line        int      `json:"line,omitempty"`
	Remediation string   `json:"remediation"`
	Explanation string   `json:"explanation,omitempty"`
}

type Report struct {
	SchemaVersion string    `json:"schemaVersion"`
	Target        string    `json:"target"`
	Compatibility string    `json:"compatibility"`
	Findings      []Finding `json:"findings"`
	Summary       Summary   `json:"summary"`
}

type Summary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Notes    int `json:"notes"`
}

type manifest struct {
	SchemaVersion any            `json:"schemaVersion"`
	ID            any            `json:"id"`
	Name          any            `json:"name"`
	Version       any            `json:"version"`
	Kinds         any            `json:"kinds"`
	EntryPoints   any            `json:"entryPoints"`
	BarWidget     map[string]any `json:"barWidget"`
	Raw           map[string]any `json:"-"`
}

var idPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
var semverPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
var stringCommandPattern = regexp.MustCompile(`(?:command\s*:|\.command\s*=)\s*["']`)
var hardCodedColorPattern = regexp.MustCompile(`#[0-9A-Fa-f]{6,8}`)

func Run(target string) Report { return RunFor(target, "omarchy-4") }

func RunFor(target, compatibility string) Report {
	abs, err := filepath.Abs(target)
	if err != nil {
		abs = target
	}
	report := Report{SchemaVersion: SchemaVersion, Target: filepath.Clean(abs), Compatibility: compatibility, Findings: []Finding{}}
	add := func(f Finding) { report.Findings = append(report.Findings, f) }
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		add(finding("OF001", "official-parity", Error, "plugin directory is not available", ".", "Pass an existing plugin directory.", "Matches the official validator's plugin-folder boundary."))
		return finish(report)
	}

	manifestPath := filepath.Join(abs, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		add(finding("OF100", "official-parity", Error, "manifest.json is missing or unreadable", "manifest.json", "Add a readable manifest.json at the plugin root.", ""))
	} else {
		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err != nil {
			add(finding("OF101", "official-parity", Error, "manifest.json is not valid JSON", "manifest.json", "Correct the JSON syntax.", err.Error()))
		} else {
			checkManifest(abs, raw, add)
		}
	}

	_ = filepath.WalkDir(abs, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		rel, _ := filepath.Rel(abs, path)
		if rel == ".git" {
			return filepath.SkipDir
		}
		if entry.Type()&os.ModeSymlink != 0 {
			add(finding("OF106", "official-parity", Error, "symlinks are not allowed inside a plugin", filepath.ToSlash(rel), "Replace the symlink with a regular file or directory.", ""))
		}
		if entry.IsDir() || !isSourceFile(path) {
			return nil
		}
		contents, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		checkSource(filepath.ToSlash(rel), string(contents), add)
		return nil
	})
	checkDocs(abs, add)
	return finish(report)
}

func checkManifest(root string, raw map[string]any, add func(Finding)) {
	if value, ok := raw["schemaVersion"].(float64); !ok || value != 1 {
		add(finding("OF102", "official-parity", Error, "schemaVersion must be the JSON number 1", "manifest.json", "Set schemaVersion to 1.", ""))
	}
	for _, field := range []string{"id", "name", "version", "kinds", "entryPoints"} {
		if _, ok := raw[field]; !ok {
			add(finding("OF103", "official-parity", Error, fmt.Sprintf("manifest is missing required field %q", field), "manifest.json", "Add the required field with the documented type.", ""))
		}
	}
	id, idOK := raw["id"].(string)
	if !idOK || id == "" || !idPattern.MatchString(id) || strings.Contains(id, "..") || strings.HasPrefix(id, "omarchy.") {
		add(finding("OF104", "official-parity", Error, "plugin id is invalid or uses the reserved omarchy.* namespace", "manifest.json", "Use a namespaced third-party ID such as dev.example.plugin.", ""))
	}
	version, versionOK := raw["version"].(string)
	if !versionOK || !semverPattern.MatchString(version) {
		add(finding("OF200", "forge", Error, "plugin version is not semantic versioning", "manifest.json", "Use a version such as 0.1.0.", "Forge publish-readiness rule."))
	}
	kindsRaw, kindsOK := raw["kinds"].([]any)
	kinds := []string{}
	if !kindsOK || len(kindsRaw) == 0 {
		add(finding("OF105", "official-parity", Error, "kinds must be a non-empty array", "manifest.json", "Declare at least one supported plugin kind.", ""))
	} else {
		known := map[string]string{"bar": "bar", "bar-widget": "barWidget", "menu": "menu", "overlay": "overlay", "panel": "panel", "service": "service"}
		for _, value := range kindsRaw {
			kind, ok := value.(string)
			if !ok {
				add(finding("OF105", "official-parity", Error, "every kind must be a string", "manifest.json", "Use supported string kind names.", ""))
				continue
			}
			kinds = append(kinds, kind)
			if _, ok := known[kind]; !ok {
				add(finding("OF107", "forge", Error, fmt.Sprintf("unknown plugin kind %q", kind), "manifest.json", "Use a kind supported by the pinned compatibility contract.", "Forge is stricter than the current official validator for unknown kinds."))
			}
		}
	}
	entries, entriesOK := raw["entryPoints"].(map[string]any)
	if barWidget, ok := raw["barWidget"].(map[string]any); ok {
		if section, present := barWidget["defaultSection"]; present {
			value, stringOK := section.(string)
			if !stringOK || (value != "left" && value != "center" && value != "right") {
				add(finding("OF112", "official-parity", Error, "barWidget.defaultSection must be left, center, or right", "manifest.json", "Use one of the three supported section values.", ""))
			}
		}
	}
	if !entriesOK {
		add(finding("OF108", "official-parity", Error, "entryPoints must be an object", "manifest.json", "Map each declared kind to a safe relative file path.", ""))
		return
	}
	known := map[string]string{"bar": "bar", "bar-widget": "barWidget", "menu": "menu", "overlay": "overlay", "panel": "panel", "service": "service"}
	for _, kind := range kinds {
		key, knownKind := known[kind]
		if !knownKind {
			continue
		}
		if _, ok := entries[key]; !ok {
			add(finding("OF109", "official-parity", Error, fmt.Sprintf("kind %q requires entryPoints.%s", kind, key), "manifest.json", "Add the required entry point.", ""))
		}
	}
	for key, value := range entries {
		path, ok := value.(string)
		if !ok || path == "" || filepath.IsAbs(path) || strings.Contains(path, "..") || strings.ContainsAny(path, "\r\n") {
			add(finding("OF110", "official-parity", Error, fmt.Sprintf("entryPoints.%s is not a safe relative path", key), "manifest.json", "Use a nonempty relative path without '..' or newlines.", ""))
			continue
		}
		info, err := os.Stat(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil || !info.Mode().IsRegular() {
			add(finding("OF111", "official-parity", Error, fmt.Sprintf("entry point file %q does not exist", path), "manifest.json", "Add the referenced regular file.", ""))
		}
	}
	if contains(kinds, "bar-widget") {
		if entry, ok := entries["barWidget"].(string); ok && entry != "" && !filepath.IsAbs(entry) && !strings.Contains(entry, "..") {
			checkBarWidget(root, entry, add)
		}
	}
}

func checkBarWidget(root, entry string, add func(Finding)) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(entry)))
	if err != nil {
		return
	}
	qml := string(data)
	checks := []struct {
		id, message, remediation string
		any                      []string
	}{
		{"OF210", "bar widget entry point has no obvious theme-token integration", "Use injected bar colors or qs.Commons Color/Style tokens.", []string{"Color.", "Style.", "barForeground", "foreground"}},
		{"OF211", "bar widget entry point has no obvious keyboard/focus handling", "Add keyboard reachability, close behavior, and a visible focus state.", []string{"PanelKeyCatcher", "Keys.on", "activeFocus"}},
		{"OF212", "bar widget entry point has no obvious loading state", "Represent asynchronous loading explicitly.", []string{"LoadingState", `state === "loading"`}},
		{"OF213", "bar widget entry point has no obvious empty state", "Represent an empty result distinctly from success.", []string{"EmptyState", `state === "empty"`}},
		{"OF214", "bar widget entry point has no obvious error state", "Expose actionable local failures without crashing the shell.", []string{"ErrorState", `state === "error"`}},
	}
	for _, check := range checks {
		matched := false
		for _, needle := range check.any {
			if strings.Contains(qml, needle) {
				matched = true
				break
			}
		}
		if !matched {
			add(finding(check.id, "forge", Warning, check.message, filepath.ToSlash(entry), check.remediation, "Heuristic entry-point source check."))
		}
	}
}

func checkDocs(root string, add func(Finding)) {
	readmePath := filepath.Join(root, "README.md")
	data, err := os.ReadFile(readmePath)
	if err != nil {
		add(finding("OF201", "forge", Error, "README.md is missing", "README.md", "Add installation, use, compatibility, privacy, update, and removal documentation.", "Publish-readiness rule."))
		return
	}
	lower := strings.ToLower(string(data))
	for _, section := range []string{"requirements", "install", "configuration", "development", "update", "removal", "privacy", "license"} {
		if !strings.Contains(lower, "## "+section) {
			add(finding("OF202", "forge", Warning, fmt.Sprintf("README is missing a %s section", section), "README.md", "Add the missing publish-readiness section.", "Heuristic Markdown heading check."))
		}
	}
	if !strings.Contains(lower, "omarchy 4") && !strings.Contains(lower, "compatib") {
		add(finding("OF203", "forge", Warning, "README does not declare Omarchy compatibility", "README.md", "State the tested Omarchy version or compatibility range.", "Heuristic documentation check."))
	}
	if _, err := os.Stat(filepath.Join(root, "LICENSE")); err != nil {
		add(finding("OF204", "forge", Error, "LICENSE is missing", "LICENSE", "Add the license file declared by the manifest.", "Publish-readiness rule."))
	}
	preview := false
	for _, name := range []string{"preview.png", "preview.jpg", "screenshot.png", "assets/preview.png", "assets/screenshot.png"} {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			preview = true
		}
	}
	if !preview && !strings.Contains(lower, "preview") && !strings.Contains(lower, "screenshot") {
		add(finding("OF205", "forge", Warning, "plugin has no preview image or explicit preview warning", "README.md", "Add a representative image or document that a preview is not yet available.", "Publish-readiness heuristic."))
	}
}

func checkSource(path, contents string, add func(Finding)) {
	lines := strings.Split(contents, "\n")
	patterns := []struct{ id, needle, message, remediation string }{
		{"OF300", "/usr/share/omarchy", "source references an Omarchy-owned write-sensitive path", "Do not write to packaged Omarchy paths; use documented user locations only."},
		{"OF301", "sudo ", "source may invoke sudo", "Remove privilege escalation from plugin runtime and installation."},
		{"OF302", "pacman ", "source may install packages", "Document requirements and let users install dependencies themselves."},
	}
	for number, line := range lines {
		for _, pattern := range patterns {
			if strings.Contains(line, pattern.needle) {
				match := finding(pattern.id, "forge", Warning, pattern.message, path, pattern.remediation, "Heuristic source-pattern match; inspect context before treating it as a defect.")
				match.Line = number + 1
				add(match)
			}
		}
		if strings.HasSuffix(strings.ToLower(path), ".qml") && stringCommandPattern.MatchString(line) {
			match := finding("OF303", "forge", Warning, "QML command appears to use a shell-like string", "", "Use array-form Process.command arguments and avoid shell interpolation.", "Heuristic source-pattern match; inspect context before treating it as unsafe command construction.")
			match.Path = path
			match.Line = number + 1
			add(match)
		}
		if strings.HasSuffix(strings.ToLower(path), ".qml") && hardCodedColorPattern.MatchString(line) {
			match := finding("OF304", "forge", Warning, "QML contains a hard-coded color", "", "Prefer Omarchy Color, Style, bar, and foreground theme tokens.", "Heuristic theme-integration pattern match.")
			match.Path = path
			match.Line = number + 1
			add(match)
		}
	}
}

func isSourceFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".qml", ".js", ".sh":
		return true
	}
	return false
}
func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func finding(id, source string, severity Severity, message, path, remediation, explanation string) Finding {
	if explanation == "" {
		if source == "official-parity" {
			explanation = "Mirrors the inspected Omarchy 4 official validator contract."
		} else {
			explanation = "Forge quality rule."
		}
	}
	return Finding{RuleID: id, Source: source, Severity: severity, Message: message, Path: path, Remediation: remediation, Explanation: explanation}
}
func finish(report Report) Report {
	sort.SliceStable(report.Findings, func(i, j int) bool {
		a, b := report.Findings[i], report.Findings[j]
		if a.RuleID != b.RuleID {
			return a.RuleID < b.RuleID
		}
		if a.Path != b.Path {
			return a.Path < b.Path
		}
		return a.Message < b.Message
	})
	for _, f := range report.Findings {
		switch f.Severity {
		case Error:
			report.Summary.Errors++
		case Warning:
			report.Summary.Warnings++
		case Note:
			report.Summary.Notes++
		}
	}
	return report
}
