package main

import (
	"fmt"
	"strings"
)

const (
	markerStart = "<!-- rules:start -->"
	markerEnd   = "<!-- rules:end -->"
)

// wrapBlock wraps the rules block between the start/end markers.
func wrapBlock(block string) string {
	return markerStart + "\n" + strings.TrimRight(block, "\n") + "\n" + markerEnd
}

// hasBlock reports whether content already contains a rules block.
func hasBlock(content string) bool {
	return strings.Contains(content, markerStart) && strings.Contains(content, markerEnd)
}

// injectBlock returns the CLAUDE.md content with the wrapped rules block
// injected. Absent markers append after existing content; present markers get
// their contents replaced (never nested).
func injectBlock(existing, block string) (string, error) {
	wrapped := wrapBlock(block)
	if existing == "" {
		return wrapped + "\n", nil
	}
	si := strings.Index(existing, markerStart)
	ei := strings.Index(existing, markerEnd)
	if si == -1 && ei == -1 {
		return strings.TrimRight(existing, "\n") + "\n\n" + wrapped + "\n", nil
	}
	if si == -1 || ei == -1 || ei < si {
		return "", fmt.Errorf("CLAUDE.md has malformed rules markers (start index %d, end index %d)", si, ei)
	}
	before := existing[:si]
	after := existing[ei+len(markerEnd):]
	return before + wrapped + after, nil
}
