package ignorer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envdiff/internal/ignorer"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envignore")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempIgnore: %v", err)
	}
	return path
}

func TestNew_Empty(t *testing.T) {
	ig := ignorer.New()
	if ig.ShouldIgnore("ANY_KEY") {
		t.Error("expected empty ignorer to not ignore any key")
	}
}

func TestAdd_AndShouldIgnore(t *testing.T) {
	ig := ignorer.New()
	ig.Add("SECRET_KEY")
	if !ig.ShouldIgnore("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be ignored")
	}
	if ig.ShouldIgnore("OTHER_KEY") {
		t.Error("expected OTHER_KEY to not be ignored")
	}
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempIgnore(t, "# comment\nSECRET_KEY\nANOTHER_KEY\n")
	ig := ignorer.New()
	if err := ig.LoadFile(path); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	for _, key := range []string{"SECRET_KEY", "ANOTHER_KEY"} {
		if !ig.ShouldIgnore(key) {
			t.Errorf("expected %q to be ignored", key)
		}
	}
}

func TestLoadFile_BlankLines(t *testing.T) {
	path := writeTempIgnore(t, "\n\nKEY_A\n\nKEY_B\n")
	ig := ignorer.New()
	if err := ig.LoadFile(path); err != nil {
		t.Fatalf("LoadFile: %v", err)
	}
	if !ig.ShouldIgnore("KEY_A") || !ig.ShouldIgnore("KEY_B") {
		t.Error("expected KEY_A and KEY_B to be ignored")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	ig := ignorer.New()
	if err := ig.LoadFile("/nonexistent/.envignore"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFilter(t *testing.T) {
	ig := ignorer.New()
	ig.Add("SECRET")
	env := map[string]string{"SECRET": "abc", "PORT": "8080", "HOST": "localhost"}
	filtered := ig.Filter(env)
	if _, ok := filtered["SECRET"]; ok {
		t.Error("expected SECRET to be filtered out")
	}
	if filtered["PORT"] != "8080" || filtered["HOST"] != "localhost" {
		t.Error("expected non-ignored keys to remain")
	}
}
