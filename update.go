package main

import (
	"errors"
	"fmt"
	"os"
)

// runUpdate refreshes managed skills and re-injects the CLAUDE.md block.
// User-modified skill files are left untouched.
func runUpdate(root string) error {
	lf, err := readLock(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no %s found in %s; run `rules init` first", lockFileName, root)
		}
		return err
	}
	skills, err := embeddedSkills()
	if err != nil {
		return err
	}
	block, err := readClaudeBlock()
	if err != nil {
		return err
	}
	out, err := updateSkills(root, skills, lf.Files)
	if err != nil {
		return err
	}
	if err := injectClaudeMd(root, block); err != nil {
		return err
	}
	lf.Version = version
	lf.InstalledAt = nowRFC3339()
	lf.Files = out.files
	if err := writeLock(root, lf); err != nil {
		return err
	}
	printUpdateSummary(out)
	return nil
}

// updateSkills refreshes each embedded skill. A file whose on-disk hash still
// matches the lock hash is rewritten; a file the user changed is skipped and
// its lock entry preserved.
func updateSkills(root string, skills []skillFile, locked map[string]string) (*installOutcome, error) {
	out := newOutcome()
	for _, sk := range skills {
		newHash := sha256Hex(sk.content)
		onDisk, exists, err := hashIfExists(skillTarget(root, sk.relPath))
		if err != nil {
			return nil, err
		}
		lockHash, inLock := locked[sk.relPath]
		switch {
		case !exists:
			err = record(out, &out.written, root, sk, newHash)
		case inLock && onDisk == lockHash:
			err = record(out, &out.updated, root, sk, newHash)
		case onDisk == newHash:
			out.unchanged = append(out.unchanged, sk.relPath)
			out.files[sk.relPath] = newHash
		case !inLock:
			out.skipped = append(out.skipped, sk.relPath)
		default:
			out.skipped = append(out.skipped, sk.relPath)
			out.files[sk.relPath] = lockHash
		}
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func printUpdateSummary(out *installOutcome) {
	fmt.Println("rules update:")
	printBucket("updated", out.updated)
	printBucket("restored", out.written)
	printBucket("unchanged", out.unchanged)
	for _, p := range out.skipped {
		fmt.Printf("  skipped %s (you modified it; not overwriting)\n", p)
	}
	fmt.Println("CLAUDE.md rules block re-injected between markers")
	fmt.Printf("summary: %d updated, %d restored, %d unchanged, %d skipped\n",
		len(out.updated), len(out.written), len(out.unchanged), len(out.skipped))
}
