// Package templates embeds and renders Forge project templates.
package templates

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"text/template"
)

//go:embed all:bar-widget/files
var templateFiles embed.FS

// Data contains validated values used by templates.
type Data struct {
	ID          string
	Name        string
	Description string
	Author      string
	Section     string
	IncludeCI   bool
}

// File is one rendered project file.
type File struct {
	Path    string
	Content []byte
	Mode    fs.FileMode
}

// Render renders the bar-widget template deterministically.
func Render(data Data) ([]File, error) {
	var names []string
	err := fs.WalkDir(templateFiles, "bar-widget/files", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !entry.IsDir() {
			names = append(names, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)

	functions := template.FuncMap{
		"jsonString": func(value string) string {
			encoded, _ := json.Marshal(value)
			return string(encoded)
		},
		"markdown": escapeMarkdown,
	}
	files := make([]File, 0, len(names))
	for _, name := range names {
		relative := strings.TrimPrefix(name, "bar-widget/files/")
		if relative == ".github/workflows/forge.yml.tmpl" && !data.IncludeCI {
			continue
		}
		content, readErr := templateFiles.ReadFile(name)
		if readErr != nil {
			return nil, readErr
		}
		parsed, parseErr := template.New(relative).Funcs(functions).Option("missingkey=error").Parse(string(content))
		if parseErr != nil {
			return nil, fmt.Errorf("parse %s: %w", relative, parseErr)
		}
		var rendered bytes.Buffer
		if executeErr := parsed.Execute(&rendered, data); executeErr != nil {
			return nil, fmt.Errorf("render %s: %w", relative, executeErr)
		}
		outputPath := strings.TrimSuffix(relative, ".tmpl")
		mode := fs.FileMode(0o644)
		if outputPath == "demo/run" || outputPath == "tests/run" {
			mode = 0o755
		}
		files = append(files, File{Path: path.Clean(outputPath), Content: rendered.Bytes(), Mode: mode})
	}
	return files, nil
}

func escapeMarkdown(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\", "`", "\\`", "*", "\\*", "_", "\\_", "{", "\\{", "}", "\\}",
		"[", "\\[", "]", "\\]", "<", "&lt;", ">", "&gt;", "#", "\\#", "+", "\\+",
		"-", "\\-", ".", "\\.", "!", "\\!", "|", "\\|",
	)
	return replacer.Replace(value)
}
