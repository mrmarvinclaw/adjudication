package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveExplicitCaseFilesExpandsGlob(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "a.txt")
	second := filepath.Join(dir, "b.txt")
	if err := os.WriteFile(first, []byte("a"), 0o644); err != nil {
		t.Fatalf("write first: %v", err)
	}
	if err := os.WriteFile(second, []byte("b"), 0o644); err != nil {
		t.Fatalf("write second: %v", err)
	}

	got, err := resolveExplicitCaseFiles([]string{filepath.Join(dir, "*.txt")})
	if err != nil {
		t.Fatalf("resolveExplicitCaseFiles returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("resolveExplicitCaseFiles returned %d files, want 2", len(got))
	}
	wantFirst, _ := filepath.Abs(first)
	wantSecond, _ := filepath.Abs(second)
	if got[0] != wantFirst || got[1] != wantSecond {
		t.Fatalf("resolveExplicitCaseFiles = %#v, want [%q %q]", got, wantFirst, wantSecond)
	}
}

func TestResolveExplicitCaseFilesRejectsProhibitedExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sign.sh")
	if err := os.WriteFile(path, []byte("echo hi\n"), 0o644); err != nil {
		t.Fatalf("write sign.sh: %v", err)
	}

	_, err := resolveExplicitCaseFiles([]string{path})
	if err == nil || !strings.Contains(err.Error(), "prohibited extension") {
		t.Fatalf("resolveExplicitCaseFiles error = %v, want prohibited extension error", err)
	}
}

func TestResolveExplicitCaseFilesRejectsUnmatchedGlob(t *testing.T) {
	dir := t.TempDir()
	_, err := resolveExplicitCaseFiles([]string{filepath.Join(dir, "*.txt")})
	if err == nil || !strings.Contains(err.Error(), "matched no files") {
		t.Fatalf("resolveExplicitCaseFiles error = %v, want unmatched glob error", err)
	}
}
