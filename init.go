package main

import (
	"fmt"
	"path/filepath"
)

// runInit installs skills and injects the CLAUDE.md rules block into root.
func runInit(root string, force bool) error {
	skills, err := embeddedSkills()
	if err != nil {
		return err
	}
	block, err := readClaudeBlock()
	if err != nil {
		return err
	}
	managed, err := managedFiles(root)
	if err != nil {
		return err
	}
	out, err := installSkills(root, skills, managed, force)
	if err != nil {
		return err
	}
	if err := injectClaudeMd(root, block); err != nil {
		return err
	}
	lf := &lockFile{Version: version, InstalledAt: nowRFC3339(), Files: out.files}
	if err := writeLock(root, lf); err != nil {
		return err
	}
	printInitSummary(out, force)
	return nil
}

// installSkills writes each skill according to the conflict rules: an existing
// file that differs from ours and is not lock-tracked is a conflict, skipped
// unless force is set.
func installSkills(root string, skills []skillFile, managed map[string]string, force bool) (*installOutcome, error) {
	out := newOutcome()
	for _, sk := range skills {
		newHash := sha256Hex(sk.content)
		target := skillTarget(root, sk.relPath)
		onDisk, exists, err := hashIfExists(target)
		if err != nil {
			return nil, err
		}
		_, tracked := managed[sk.relPath]
		switch {
		case !exists:
			err = record(out, &out.written, root, sk, newHash)
		case onDisk == newHash:
			out.unchanged = append(out.unchanged, sk.relPath)
			out.files[sk.relPath] = newHash
		case tracked || force:
			err = record(out, &out.updated, root, sk, newHash)
		default:
			out.skipped = append(out.skipped, sk.relPath)
		}
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func skillTarget(root, relPath string) string {
	return filepath.Join(root, filepath.FromSlash(relPath))
}

// record writes a skill file, appends its path to bucket, and records its hash.
func record(out *installOutcome, bucket *[]string, root string, sk skillFile, hash string) error {
	if err := writeManagedFile(root, sk.relPath, sk.content); err != nil {
		return err
	}
	*bucket = append(*bucket, sk.relPath)
	out.files[sk.relPath] = hash
	return nil
}

func printInitSummary(out *installOutcome, force bool) {
	fmt.Println("rules init:")
	printBucket("written", out.written)
	printBucket("updated", out.updated)
	printBucket("unchanged", out.unchanged)
	printBucket("skipped(conflict)", out.skipped)
	fmt.Printf("CLAUDE.md rules block injected between %s and %s\n", markerStart, markerEnd)
	fmt.Printf("summary: %d written, %d updated, %d unchanged, %d skipped(conflict)\n",
		len(out.written), len(out.updated), len(out.unchanged), len(out.skipped))
	if len(out.skipped) > 0 && !force {
		fmt.Println("note: re-run `rules init --force` to overwrite conflicting files")
	}
}

func printBucket(label string, items []string) {
	for _, p := range items {
		fmt.Printf("  %s %s\n", label, p)
	}
}
