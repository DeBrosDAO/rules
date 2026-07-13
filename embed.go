package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed all:content
var contentFS embed.FS

const (
	claudeBlockPath = "content/claude-block.md"
	skillsDir       = "content/skills"
)

// readClaudeBlock returns the embedded hard-rules block that gets injected into
// a project's CLAUDE.md.
func readClaudeBlock() (string, error) {
	b, err := contentFS.ReadFile(claudeBlockPath)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded %s: %w", claudeBlockPath, err)
	}
	return string(b), nil
}

// skillFile is one embedded skill and the project-relative path it installs to.
type skillFile struct {
	name    string // skill name, e.g. "rules-bugboard"
	relPath string // canonical (forward-slash) install path under the project root
	content []byte
}

// embeddedSkills walks the embedded skills directory and returns every
// SKILL.md, namespaced under .claude/skills/<name>/SKILL.md.
func embeddedSkills() ([]skillFile, error) {
	entries, err := fs.ReadDir(contentFS, skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded %s: %w", skillsDir, err)
	}
	var skills []skillFile
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		embedPath := skillsDir + "/" + name + "/SKILL.md"
		b, err := contentFS.ReadFile(embedPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded skill %s: %w", name, err)
		}
		rel := filepath.ToSlash(filepath.Join(".claude", "skills", name, "SKILL.md"))
		skills = append(skills, skillFile{name: name, relPath: rel, content: b})
	}
	if len(skills) == 0 {
		return nil, fmt.Errorf("no embedded skills found under %s", skillsDir)
	}
	return skills, nil
}
