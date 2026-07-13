package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// runDoctor reports the health of a rules install and returns false on any
// drift, missing file, or absent CLAUDE.md block.
func runDoctor(root string) (bool, error) {
	lf, err := readLock(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("no %s found in %s; run `rules init`\n", lockFileName, root)
			return false, nil
		}
		return false, err
	}
	fmt.Printf("rules doctor (%s):\n", root)
	fmt.Printf("  cli version:       %s\n", version)
	fmt.Printf("  installed version: %s\n", lf.Version)
	ok := true
	if lf.Version != version {
		fmt.Println("  drift: installed version differs from CLI (run `rules update`)")
		ok = false
	}
	if !checkFiles(root, lf.Files) {
		ok = false
	}
	if !checkClaudeMd(root) {
		ok = false
	}
	if ok {
		fmt.Println("all good")
	}
	return ok, nil
}

// checkFiles verifies every lock-tracked file exists and matches its hash.
func checkFiles(root string, files map[string]string) bool {
	names := make([]string, 0, len(files))
	for k := range files {
		names = append(names, k)
	}
	sort.Strings(names)
	ok := true
	for _, rel := range names {
		target := filepath.Join(root, filepath.FromSlash(rel))
		onDisk, exists, err := hashIfExists(target)
		switch {
		case err != nil:
			fmt.Printf("  skill %s: error (%v)\n", rel, err)
			ok = false
		case !exists:
			fmt.Printf("  skill %s: MISSING\n", rel)
			ok = false
		case onDisk != files[rel]:
			fmt.Printf("  skill %s: MODIFIED\n", rel)
			ok = false
		default:
			fmt.Printf("  skill %s: ok\n", rel)
		}
	}
	return ok
}

func checkClaudeMd(root string) bool {
	b, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		fmt.Println("  CLAUDE.md: MISSING")
		return false
	}
	if !hasBlock(string(b)) {
		fmt.Println("  CLAUDE.md: rules block MISSING")
		return false
	}
	fmt.Println("  CLAUDE.md: rules block present")
	return true
}
