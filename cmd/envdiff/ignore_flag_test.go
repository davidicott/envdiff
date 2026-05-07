package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempIgnoreFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envignore")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempIgnoreFile: %v", err)
	}
	return path
}

func TestLoadIgnorer_NoArgs(t *testing.T) {
	ig, err := loadIgnorer("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ig.ShouldIgnore("ANYTHING") {
		t.Error("expected no-op ignorer")
	}
}

func TestLoadIgnorer_FromFile(t *testing.T) {
	path := writeTempIgnoreFile(t, "SECRET\nTOKEN\n")
	ig, err := loadIgnorer(path, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ig.ShouldIgnore("SECRET") || !ig.ShouldIgnore("TOKEN") {
		t.Error("expected SECRET and TOKEN to be ignored")
	}
}

func TestLoadIgnorer_InlineKeys(t *testing.T) {
	ig, err := loadIgnorer("", []string{"MY_KEY", "OTHER_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ig.ShouldIgnore("MY_KEY") || !ig.ShouldIgnore("OTHER_KEY") {
		t.Error("expected inline keys to be ignored")
	}
}

func TestLoadIgnorer_FileNotFound(t *testing.T) {
	_, err := loadIgnorer("/nonexistent/.envignore", nil)
	if err == nil {
		t.Error("expected error for missing ignore file")
	}
}

func TestLoadIgnorer_FileAndInlineKeys(t *testing.T) {
	path := writeTempIgnoreFile(t, "FILE_KEY\n")
	ig, err := loadIgnorer(path, []string{"INLINE_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ig.ShouldIgnore("FILE_KEY") || !ig.ShouldIgnore("INLINE_KEY") {
		t.Error("expected both file and inline keys to be ignored")
	}
}
