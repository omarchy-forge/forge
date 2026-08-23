package scaffold

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	maximumTextReferenceSize  = 256 * 1024
	maximumImageReferenceSize = 5 * 1024 * 1024
)

// Reference is a validated local project brief or image copied into a scaffold.
type Reference struct {
	SourcePath  string
	ProjectPath string
	Kind        string
	Size        int
	SHA256      string
	content     []byte
}

// CreateReferenceDirectory validates the target and creates the guided input directory.
func CreateReferenceDirectory(options Options) (string, error) {
	target, err := validate(options)
	if err != nil {
		return "", err
	}
	directory := filepath.Join(target, "references")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return "", fmt.Errorf("create reference directory: %w", err)
	}
	return directory, nil
}

// PrepareReferencesInDirectory validates every regular file in a guided directory.
func PrepareReferencesInDirectory(directory string) ([]Reference, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("read reference directory: %w", err)
	}
	if len(entries) > 10 {
		return nil, fmt.Errorf("at most 10 reference files are allowed")
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.Name() == ".gitkeep" {
			continue
		}
		paths = append(paths, filepath.Join(directory, entry.Name()))
	}
	references, err := PrepareReferences(paths)
	if err != nil {
		return nil, err
	}
	for index := range references {
		references[index].ProjectPath = filepath.ToSlash(filepath.Join("references", filepath.Base(references[index].SourcePath)))
	}
	return references, nil
}

// PrepareReferences validates bounded regular files without executing or decoding them.
func PrepareReferences(paths []string) ([]Reference, error) {
	if len(paths) > 10 {
		return nil, fmt.Errorf("at most 10 reference files are allowed")
	}
	references := make([]Reference, 0, len(paths))
	seen := map[string]bool{}
	for index, source := range paths {
		absolute, err := filepath.Abs(source)
		if err != nil {
			return nil, fmt.Errorf("resolve reference %q: %w", source, err)
		}
		absolute = filepath.Clean(absolute)
		if seen[absolute] {
			return nil, fmt.Errorf("reference file was selected more than once: %s", absolute)
		}
		seen[absolute] = true
		info, err := os.Lstat(absolute)
		if err != nil {
			return nil, fmt.Errorf("inspect reference %s: %w", absolute, err)
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("reference must be a regular file, not a symlink or directory: %s", absolute)
		}
		extension := strings.ToLower(filepath.Ext(absolute))
		kind, maximum, err := referenceType(extension)
		if err != nil {
			return nil, fmt.Errorf("reference %s: %w", absolute, err)
		}
		if info.Size() > int64(maximum) {
			return nil, fmt.Errorf("reference %s exceeds the %d-byte %s limit", absolute, maximum, kind)
		}
		content, err := os.ReadFile(absolute)
		if err != nil {
			return nil, fmt.Errorf("read reference %s: %w", absolute, err)
		}
		if kind == "text" {
			if !utf8.Valid(content) || bytes.IndexByte(content, 0) >= 0 {
				return nil, fmt.Errorf("text reference must be valid UTF-8 without NUL bytes: %s", absolute)
			}
		} else if extension == ".svg" {
			if !validSVG(content) {
				return nil, fmt.Errorf("SVG reference must be valid UTF-8 XML with an svg root element: %s", absolute)
			}
		} else if !validImageSignature(content, extension) {
			return nil, fmt.Errorf("image signature does not match %s: %s", extension, absolute)
		}
		digest := sha256.Sum256(content)
		projectExtension := extension
		if projectExtension == ".jpeg" {
			projectExtension = ".jpg"
		}
		references = append(references, Reference{
			SourcePath:  absolute,
			ProjectPath: fmt.Sprintf("references/reference-%02d%s", index+1, projectExtension),
			Kind:        kind,
			Size:        len(content),
			SHA256:      hex.EncodeToString(digest[:]),
			content:     content,
		})
	}
	return references, nil
}

func validatePreparedReferences(references []Reference) error {
	if len(references) > 10 {
		return fmt.Errorf("at most 10 reference files are allowed")
	}
	for index, reference := range references {
		extension := strings.ToLower(filepath.Ext(reference.ProjectPath))
		kind, maximum, err := referenceType(extension)
		if err != nil {
			return fmt.Errorf("invalid prepared reference %d: %w", index+1, err)
		}
		deterministic := fmt.Sprintf("references/reference-%02d%s", index+1, extension)
		inPlace := filepath.ToSlash(filepath.Join("references", filepath.Base(reference.SourcePath)))
		if (reference.ProjectPath != deterministic && reference.ProjectPath != inPlace) || reference.Kind != kind || reference.Size != len(reference.content) || len(reference.content) > maximum {
			return fmt.Errorf("invalid prepared reference metadata for %s", reference.ProjectPath)
		}
		digest := sha256.Sum256(reference.content)
		if reference.SHA256 != hex.EncodeToString(digest[:]) {
			return fmt.Errorf("invalid prepared reference digest for %s", reference.ProjectPath)
		}
	}
	return nil
}

func referenceType(extension string) (string, int, error) {
	switch extension {
	case ".txt", ".md":
		return "text", maximumTextReferenceSize, nil
	case ".png", ".jpg", ".jpeg", ".webp", ".svg":
		return "image", maximumImageReferenceSize, nil
	default:
		return "", 0, fmt.Errorf("supported extensions are .txt, .md, .png, .jpg, .jpeg, .webp, and .svg")
	}
}

// validSVG parses XML only to establish the file type. Forge never renders or
// executes reference SVG, and encoding/xml does not resolve external entities.
func validSVG(content []byte) bool {
	if !utf8.Valid(content) || bytes.IndexByte(content, 0) >= 0 {
		return false
	}
	decoder := xml.NewDecoder(bytes.NewReader(content))
	foundRoot := false
	for {
		token, err := decoder.Token()
		if err != nil {
			return err == io.EOF && foundRoot
		}
		if start, ok := token.(xml.StartElement); ok {
			if !foundRoot && strings.ToLower(start.Name.Local) != "svg" {
				return false
			}
			foundRoot = true
		}
	}
}

func validImageSignature(content []byte, extension string) bool {
	switch extension {
	case ".png":
		return len(content) >= 8 && bytes.Equal(content[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10})
	case ".jpg", ".jpeg":
		return len(content) >= 3 && content[0] == 0xff && content[1] == 0xd8 && content[2] == 0xff
	case ".webp":
		return len(content) >= 12 && string(content[:4]) == "RIFF" && string(content[8:12]) == "WEBP"
	default:
		return false
	}
}
