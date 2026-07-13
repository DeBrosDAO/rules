package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const lockFileName = ".rules.lock"

// lockFile is the on-disk manifest of every file rules manages in a project.
type lockFile struct {
	Version     string            `json:"version"`
	InstalledAt string            `json:"installed_at"`
	Files       map[string]string `json:"files"`
}

func lockPath(root string) string {
	return filepath.Join(root, lockFileName)
}

// readLock loads the lock file. A missing file returns os.ErrNotExist.
func readLock(root string) (*lockFile, error) {
	b, err := os.ReadFile(lockPath(root))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("failed to read %s: %w", lockFileName, err)
	}
	var lf lockFile
	if err := json.Unmarshal(b, &lf); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", lockFileName, err)
	}
	if lf.Files == nil {
		lf.Files = map[string]string{}
	}
	return &lf, nil
}

// writeLock persists the lock file with deterministic (sorted) key ordering.
func writeLock(root string, lf *lockFile) error {
	b, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", lockFileName, err)
	}
	b = append(b, '\n')
	if err := os.WriteFile(lockPath(root), b, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", lockFileName, err)
	}
	return nil
}

// managedFiles returns the files map from an existing lock, or an empty map
// when no lock is present.
func managedFiles(root string) (map[string]string, error) {
	lf, err := readLock(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	return lf.Files, nil
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// hashIfExists returns the sha256 of a file, whether it exists, and any error
// other than not-exist.
func hashIfExists(path string) (hash string, exists bool, err error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to read %s: %w", path, err)
	}
	return sha256Hex(b), true, nil
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
