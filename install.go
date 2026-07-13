package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// installOutcome accumulates per-file results across an init or update run.
type installOutcome struct {
	written   []string
	updated   []string
	unchanged []string
	skipped   []string
	files     map[string]string
}

func newOutcome() *installOutcome {
	return &installOutcome{files: map[string]string{}}
}

// writeManagedFile writes content to a project-relative path, creating parent
// directories as needed.
func writeManagedFile(root, relPath string, content []byte) error {
	target := filepath.Join(root, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", relPath, err)
	}
	if err := os.WriteFile(target, content, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", relPath, err)
	}
	return nil
}

// injectClaudeMd injects (or refreshes) the rules block in the project's
// CLAUDE.md. It only ever touches content between the markers.
func injectClaudeMd(root, block string) error {
	path := filepath.Join(root, "CLAUDE.md")
	existing := ""
	if b, err := os.ReadFile(path); err == nil {
		existing = string(b)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read CLAUDE.md: %w", err)
	}
	next, err := injectBlock(existing, block)
	if err != nil {
		return err
	}
	if next == existing {
		return nil
	}
	if err := os.WriteFile(path, []byte(next), 0o644); err != nil {
		return fmt.Errorf("failed to write CLAUDE.md: %w", err)
	}
	return nil
}
