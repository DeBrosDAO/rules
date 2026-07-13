package main

import (
	"strings"
	"testing"
)

func TestInjectBlock_replacesNotNests(t *testing.T) {
	block := "RULES"
	first, err := injectBlock("", block)
	if err != nil {
		t.Fatal(err)
	}
	second, err := injectBlock(first, block)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(second, markerStart) != 1 {
		t.Errorf("markers nested: %q", second)
	}
	if first != second {
		t.Errorf("re-injection not idempotent:\n%q\n%q", first, second)
	}
}

func TestInjectBlock_preservesSurroundingContent(t *testing.T) {
	existing := "before\n" + wrapBlock("OLD") + "\nafter\n"
	next, err := injectBlock(existing, "NEW")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(next, "before") || !strings.Contains(next, "after") {
		t.Errorf("surrounding content lost: %q", next)
	}
	if strings.Contains(next, "OLD") {
		t.Errorf("old block content not replaced: %q", next)
	}
	if !strings.Contains(next, "NEW") {
		t.Errorf("new block content missing: %q", next)
	}
}

func TestInjectBlock_malformedMarkers(t *testing.T) {
	// end before start is malformed and must error rather than corrupt the file.
	existing := markerEnd + "\nstuff\n" + markerStart + "\n"
	if _, err := injectBlock(existing, "X"); err == nil {
		t.Fatal("expected error for malformed markers")
	}
}
