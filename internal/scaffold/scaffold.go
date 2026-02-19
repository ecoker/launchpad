package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ehrencoker/agent-kit/templates"
)

// Options configures a scaffold run.
type Options struct {
	TargetDir string
	ProfileID string
	AddonIDs  []string
	Force     bool
}

// Result is returned after a successful scaffold.
type Result struct {
	Profile      *Profile
	OutputPath   string
	CreatedFiles []string
}

// Run executes the scaffold: copies core + profile + addons into the target directory.
func Run(opts Options) (*Result, error) {
	profile := FindProfile(opts.ProfileID)
	if profile == nil {
		return nil, fmt.Errorf("unknown profile %q — run 'launchpad list' to see options", opts.ProfileID)
	}

	outputPath, err := filepath.Abs(opts.TargetDir)
	if err != nil {
		return nil, fmt.Errorf("resolving target path: %w", err)
	}

	// Safety check
	if !opts.Force {
		entries, _ := os.ReadDir(outputPath)
		if len(entries) > 0 {
			return nil, fmt.Errorf("directory %s is not empty — re-run with --force to overwrite", outputPath)
		}
	}

	if err := os.MkdirAll(outputPath, 0o755); err != nil {
		return nil, fmt.Errorf("creating target directory: %w", err)
	}

	projectName := filepath.Base(outputPath)
	replacements := map[string]string{
		"{{PROJECT_NAME}}": projectName,
	}

	// 1. Core templates (always)
	if err := copyEmbeddedDir("core", outputPath, replacements, opts.Force); err != nil {
		return nil, fmt.Errorf("copying core templates: %w", err)
	}

	// 2. Profile templates
	if err := copyEmbeddedDir("profiles/"+profile.Dir, outputPath, replacements, opts.Force); err != nil {
		return nil, fmt.Errorf("copying profile templates: %w", err)
	}

	// 3. Add-on templates
	for _, addonID := range opts.AddonIDs {
		addon := FindAddon(addonID)
		if addon == nil {
			continue
		}
		if err := copyEmbeddedDir("addons/"+addon.Dir, outputPath, replacements, opts.Force); err != nil {
			return nil, fmt.Errorf("copying addon %q templates: %w", addonID, err)
		}
	}

	// Collect created files
	var created []string
	_ = filepath.WalkDir(outputPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		created = append(created, path)
		return nil
	})

	return &Result{
		Profile:      profile,
		OutputPath:   outputPath,
		CreatedFiles: created,
	}, nil
}

// copyEmbeddedDir copies files from the embedded FS into the target directory.
func copyEmbeddedDir(srcDir, destDir string, replacements map[string]string, overwrite bool) error {
	return fs.WalkDir(templates.FS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path by stripping the source prefix
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		// Check overwrite
		if !overwrite {
			if _, statErr := os.Stat(destPath); statErr == nil {
				return fmt.Errorf("refusing to overwrite %s", destPath)
			}
		}

		// Read from embedded FS
		data, err := templates.FS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		// Apply replacements for text files
		content := data
		if isTextFile(path) {
			text := string(data)
			for token, value := range replacements {
				text = strings.ReplaceAll(text, token, value)
			}
			content = []byte(text)
		}

		// Ensure parent dir exists
		if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
			return err
		}

		return os.WriteFile(destPath, content, 0o644)
	})
}

func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".txt", ".json", ".yaml", ".yml", ".toml", ".env",
		".ts", ".js", ".go", ".py", ".ex", ".exs", ".cs", ".php",
		".html", ".css", ".scss", ".sql":
		return true
	case "":
		return true // files without extension (Makefile, Dockerfile, etc.)
	default:
		return false
	}
}
