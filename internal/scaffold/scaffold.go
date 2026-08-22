// Package scaffold safely creates deterministic plugin projects.
package scaffold

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	projecttemplates "github.com/omarchy-forge/forge/templates"
)

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	idPattern   = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

// Options controls project generation.
type Options struct {
	Directory   string
	ID          string
	Name        string
	Description string
	Author      string
	License     string
	Section     string
	IncludeCI   bool
	AgentReady  bool
	InitGit     bool
	DryRun      bool
	Force       bool
	HomeDir     string
	OmarchyPath string
}

// Change describes one planned filesystem change.
type Change struct {
	Action string
	Path   string
}

// Result describes a completed or dry-run generation.
type Result struct {
	Directory string
	Changes   []Change
}

// DisplayName derives a human-readable name from an output directory.
func DisplayName(directory string) string {
	base := filepath.Base(filepath.Clean(directory))
	parts := strings.FieldsFunc(base, func(r rune) bool { return r == '-' || r == '_' })
	for i, part := range parts {
		if part != "" {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// Generate validates, plans, and optionally writes a plugin project.
func Generate(options Options) (Result, error) {
	target, err := validate(options)
	if err != nil {
		return Result{}, err
	}

	files, err := projecttemplates.Render(projecttemplates.Data{
		ID:          options.ID,
		Name:        options.Name,
		Description: options.Description,
		Author:      options.Author,
		Section:     options.Section,
		IncludeCI:   options.IncludeCI,
		AgentReady:  options.AgentReady,
	})
	if err != nil {
		return Result{}, fmt.Errorf("render template: %w", err)
	}

	result := Result{Directory: target}
	for _, file := range files {
		destination := filepath.Join(target, filepath.FromSlash(file.Path))
		action := "CREATE"
		if info, statErr := os.Lstat(destination); statErr == nil {
			if !info.Mode().IsRegular() {
				return Result{}, fmt.Errorf("generated file collision is not a regular file: %s", destination)
			}
			action = "OVERWRITE"
		} else if !errors.Is(statErr, os.ErrNotExist) {
			return Result{}, fmt.Errorf("inspect %s: %w", destination, statErr)
		}
		result.Changes = append(result.Changes, Change{Action: action, Path: destination})
	}

	if options.DryRun {
		return result, nil
	}
	for i, file := range files {
		destination := result.Changes[i].Path
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return Result{}, fmt.Errorf("create directory for %s: %w", destination, err)
		}
		if err := os.WriteFile(destination, file.Content, file.Mode); err != nil {
			return Result{}, fmt.Errorf("write %s: %w", destination, err)
		}
	}
	if options.InitGit {
		command := exec.Command("git", "init", "--initial-branch=main", target)
		if output, commandErr := command.CombinedOutput(); commandErr != nil {
			return Result{}, fmt.Errorf("initialize Git repository: %w: %s", commandErr, strings.TrimSpace(string(output)))
		}
	}
	return result, nil
}

func validate(options Options) (string, error) {
	if strings.TrimSpace(options.Directory) == "" {
		return "", errors.New("output directory is required")
	}
	target, err := filepath.Abs(options.Directory)
	if err != nil {
		return "", fmt.Errorf("resolve output directory: %w", err)
	}
	target = filepath.Clean(target)
	if err := rejectSymlinkAncestors(target); err != nil {
		return "", err
	}
	for label, root := range map[string]string{"filesystem root": string(filepath.Separator), "home directory": options.HomeDir, "Omarchy installation": options.OmarchyPath} {
		if root == "" {
			continue
		}
		absoluteRoot, rootErr := filepath.Abs(root)
		if rootErr == nil {
			cleanRoot := filepath.Clean(absoluteRoot)
			if target == cleanRoot || (label == "Omarchy installation" && pathWithin(target, cleanRoot)) {
				return "", fmt.Errorf("refusing to generate into the %s: %s", label, target)
			}
		}
	}
	if info, statErr := os.Lstat(target); statErr == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("output directory may not be a symlink: %s", target)
		}
		if !info.IsDir() {
			return "", fmt.Errorf("output path is not a directory: %s", target)
		}
		entries, readErr := os.ReadDir(target)
		if readErr != nil {
			return "", fmt.Errorf("inspect output directory: %w", readErr)
		}
		if len(entries) > 0 && !onlyGitDirectory(entries) && !options.Force {
			return "", fmt.Errorf("output directory is not empty; inspect it and pass --force to overwrite generated file collisions: %s", target)
		}
		if options.Force {
			walkErr := filepath.WalkDir(target, func(name string, entry os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if entry.Type()&os.ModeSymlink != 0 {
					return fmt.Errorf("output directory contains a symlink: %s", name)
				}
				if entry.IsDir() && entry.Name() == ".git" {
					return filepath.SkipDir
				}
				return nil
			})
			if walkErr != nil {
				return "", walkErr
			}
		}
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return "", fmt.Errorf("inspect output directory: %w", statErr)
	}

	slug := filepath.Base(target)
	if !slugPattern.MatchString(slug) || strings.Contains(slug, "--") {
		return "", fmt.Errorf("directory name %q must use lowercase letters, digits, and single hyphens", slug)
	}
	if !idPattern.MatchString(options.ID) || !strings.Contains(options.ID, ".") || strings.Contains(options.ID, "..") || strings.HasPrefix(options.ID, "omarchy.") {
		return "", fmt.Errorf("plugin ID %q must be namespaced, use letters, digits, '.', '_', or '-', and not use omarchy.*", options.ID)
	}
	for label, constraint := range map[string]struct {
		value string
		max   int
	}{"name": {options.Name, 80}, "description": {options.Description, 240}, "author": {options.Author, 120}} {
		if strings.TrimSpace(constraint.value) == "" {
			return "", fmt.Errorf("%s is required", label)
		}
		if len(constraint.value) > constraint.max || strings.ContainsAny(constraint.value, "\r\n\x00") {
			return "", fmt.Errorf("%s must be a single line of at most %d bytes", label, constraint.max)
		}
	}
	if options.License != "MIT" {
		return "", fmt.Errorf("unsupported license %q; available license: MIT", options.License)
	}
	if options.Section != "left" && options.Section != "center" && options.Section != "right" {
		return "", fmt.Errorf("section %q must be left, center, or right", options.Section)
	}
	return target, nil
}

func onlyGitDirectory(entries []os.DirEntry) bool {
	return len(entries) == 1 && entries[0].Name() == ".git" && entries[0].IsDir()
}

func rejectSymlinkAncestors(target string) error {
	for current := target; ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output path contains a symlink: %s", current)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect output path: %w", err)
		}
		parent := filepath.Dir(current)
		if parent == current {
			return nil
		}
	}
}

func pathWithin(candidate, root string) bool {
	relative, err := filepath.Rel(root, candidate)
	return err == nil && relative != "." && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

// SortedPaths returns deterministic slash-separated paths for tests and reports.
func SortedPaths(changes []Change) []string {
	paths := make([]string, 0, len(changes))
	for _, change := range changes {
		paths = append(paths, filepath.ToSlash(change.Path))
	}
	sort.Strings(paths)
	return paths
}
