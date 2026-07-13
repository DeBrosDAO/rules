package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func skillRel(name string) string {
	return filepath.Join(".claude", "skills", name, "SKILL.md")
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func embeddedSkillContent(t *testing.T, name string) []byte {
	t.Helper()
	skills, err := embeddedSkills()
	if err != nil {
		t.Fatalf("embeddedSkills: %v", err)
	}
	for _, sk := range skills {
		if sk.name == name {
			return sk.content
		}
	}
	t.Fatalf("skill %s not embedded", name)
	return nil
}

func TestInit_emptyDir(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir, false); err != nil {
		t.Fatalf("runInit: %v", err)
	}

	skills, err := embeddedSkills()
	if err != nil {
		t.Fatalf("embeddedSkills: %v", err)
	}
	if len(skills) != 3 {
		t.Fatalf("expected 3 embedded skills, got %d", len(skills))
	}
	for _, sk := range skills {
		got := readFile(t, filepath.Join(dir, filepath.FromSlash(sk.relPath)))
		if got != string(sk.content) {
			t.Errorf("skill %s content mismatch", sk.name)
		}
	}

	claude := readFile(t, filepath.Join(dir, "CLAUDE.md"))
	if !hasBlock(claude) {
		t.Fatalf("CLAUDE.md missing rules block: %q", claude)
	}
	if !strings.Contains(claude, "DeBros Engineering Rules") {
		t.Errorf("CLAUDE.md missing block body")
	}

	lf, err := readLock(dir)
	if err != nil {
		t.Fatalf("readLock: %v", err)
	}
	if lf.Version != version {
		t.Errorf("lock version = %q, want %q", lf.Version, version)
	}
	if len(lf.Files) != 3 {
		t.Errorf("lock files = %d, want 3", len(lf.Files))
	}
	for _, sk := range skills {
		if lf.Files[sk.relPath] != sha256Hex(sk.content) {
			t.Errorf("lock hash mismatch for %s", sk.relPath)
		}
	}
	if _, ok := lf.Files["CLAUDE.md"]; ok {
		t.Errorf("CLAUDE.md must not be hash-tracked in lock")
	}
}

func TestInit_existingClaudeNoMarkers(t *testing.T) {
	dir := t.TempDir()
	original := "# My Project\n\nSome existing rules here.\n"
	if err := os.WriteFile(filepath.Join(dir, "CLAUDE.md"), []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runInit(dir, false); err != nil {
		t.Fatalf("runInit: %v", err)
	}
	claude := readFile(t, filepath.Join(dir, "CLAUDE.md"))
	if !strings.Contains(claude, "# My Project") {
		t.Errorf("original content not preserved")
	}
	if !strings.Contains(claude, "Some existing rules here.") {
		t.Errorf("original body not preserved")
	}
	if !hasBlock(claude) {
		t.Errorf("rules block not appended")
	}
	if strings.Index(claude, "# My Project") > strings.Index(claude, markerStart) {
		t.Errorf("block appended before original content")
	}
}

func TestInit_idempotent(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir, false); err != nil {
		t.Fatalf("first init: %v", err)
	}
	if err := runInit(dir, false); err != nil {
		t.Fatalf("second init: %v", err)
	}
	claude := readFile(t, filepath.Join(dir, "CLAUDE.md"))
	if n := strings.Count(claude, markerStart); n != 1 {
		t.Errorf("expected 1 start marker, got %d", n)
	}
	if n := strings.Count(claude, markerEnd); n != 1 {
		t.Errorf("expected 1 end marker, got %d", n)
	}
}

func TestUpdate_userEditedSkillSkipped(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir, false); err != nil {
		t.Fatalf("init: %v", err)
	}
	edited := filepath.Join(dir, filepath.FromSlash(skillRel("rules-orama")))
	if err := os.WriteFile(edited, []byte("# my own edits\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runUpdate(dir); err != nil {
		t.Fatalf("update: %v", err)
	}
	if got := readFile(t, edited); got != "# my own edits\n" {
		t.Errorf("edited skill was overwritten: %q", got)
	}
	// An untouched skill should still be present and correct.
	other := filepath.Join(dir, filepath.FromSlash(skillRel("rules-bugboard")))
	if got := readFile(t, other); got != string(embeddedSkillContent(t, "rules-bugboard")) {
		t.Errorf("untouched skill not refreshed")
	}
	// Lock keeps the edited file's original hash so it stays managed.
	lf, err := readLock(dir)
	if err != nil {
		t.Fatal(err)
	}
	if lf.Files[skillRel("rules-orama")] != sha256Hex(embeddedSkillContent(t, "rules-orama")) {
		t.Errorf("lock hash for edited skill should remain the original embedded hash")
	}
}

func TestUpdate_noEditsRefreshes(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir, false); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := runUpdate(dir); err != nil {
		t.Fatalf("update: %v", err)
	}
	for _, name := range []string{"rules-bugboard", "rules-orama", "rules-rootwallet"} {
		p := filepath.Join(dir, filepath.FromSlash(skillRel(name)))
		if got := readFile(t, p); got != string(embeddedSkillContent(t, name)) {
			t.Errorf("skill %s not matching embedded after update", name)
		}
	}
}

func TestUpdate_requiresLock(t *testing.T) {
	dir := t.TempDir()
	err := runUpdate(dir)
	if err == nil {
		t.Fatal("expected error when lock is absent")
	}
	if !strings.Contains(err.Error(), "rules init") {
		t.Errorf("error should suggest `rules init`, got: %v", err)
	}
}

func TestInit_conflictWithoutForce(t *testing.T) {
	dir := t.TempDir()
	// Pre-existing untracked skill file with different content.
	target := filepath.Join(dir, filepath.FromSlash(skillRel("rules-orama")))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("user's own skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runInit(dir, false); err != nil {
		t.Fatalf("init: %v", err)
	}
	if got := readFile(t, target); got != "user's own skill\n" {
		t.Errorf("conflicting file overwritten without --force: %q", got)
	}
	lf, err := readLock(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := lf.Files[skillRel("rules-orama")]; ok {
		t.Errorf("skipped conflict should not be recorded in lock")
	}
	// Other skills should still install.
	if _, ok := lf.Files[skillRel("rules-bugboard")]; !ok {
		t.Errorf("non-conflicting skill not installed")
	}
}

func TestInit_conflictWithForce(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, filepath.FromSlash(skillRel("rules-orama")))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("user's own skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := runInit(dir, true); err != nil {
		t.Fatalf("init --force: %v", err)
	}
	if got := readFile(t, target); got != string(embeddedSkillContent(t, "rules-orama")) {
		t.Errorf("--force did not overwrite conflicting file")
	}
	lf, err := readLock(dir)
	if err != nil {
		t.Fatal(err)
	}
	if lf.Files[skillRel("rules-orama")] != sha256Hex(embeddedSkillContent(t, "rules-orama")) {
		t.Errorf("forced file not recorded in lock")
	}
}

func TestDoctor_cleanAndDrift(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir, false); err != nil {
		t.Fatalf("init: %v", err)
	}
	ok, err := runDoctor(dir)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if !ok {
		t.Errorf("doctor should report clean right after init")
	}
	// Modify a managed file -> doctor detects drift.
	target := filepath.Join(dir, filepath.FromSlash(skillRel("rules-orama")))
	if err := os.WriteFile(target, []byte("tampered\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	ok, err = runDoctor(dir)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if ok {
		t.Errorf("doctor should report drift after modification")
	}
}

func TestDoctor_noLock(t *testing.T) {
	dir := t.TempDir()
	ok, err := runDoctor(dir)
	if err != nil {
		t.Fatalf("doctor: %v", err)
	}
	if ok {
		t.Errorf("doctor should report not-ok without a lock file")
	}
}
