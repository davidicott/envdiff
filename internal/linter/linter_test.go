package linter_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/linter"
)

func findWarning(warnings []linter.Warning, key, fragment string) bool {
	for _, w := range warnings {
		if w.Key == key && strings.Contains(w.Message, fragment) {
			return true
		}
	}
	return false
}

func TestLint_CleanMap(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	}
	warnings := linter.Lint(env)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"database_url": "postgres://localhost/db",
	}
	warnings := linter.Lint(env)
	if !findWarning(warnings, "database_url", "lowercase") {
		t.Errorf("expected lowercase warning for 'database_url', got: %v", warnings)
	}
}

func TestLint_UnresolvedPlaceholder(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "${SECRET_KEY}",
	}
	warnings := linter.Lint(env)
	if !findWarning(warnings, "SECRET_KEY", "unresolved placeholder") {
		t.Errorf("expected placeholder warning for 'SECRET_KEY', got: %v", warnings)
	}
}

func TestLint_ValueTooLong(t *testing.T) {
	env := map[string]string{
		"BIG_VALUE": strings.Repeat("x", 1025),
	}
	warnings := linter.Lint(env)
	if !findWarning(warnings, "BIG_VALUE", "exceeds") {
		t.Errorf("expected length warning for 'BIG_VALUE', got: %v", warnings)
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	env := map[string]string{
		"mixed_Key": "${UNSET}",
		"GOOD_KEY":  "good_value",
	}
	warnings := linter.Lint(env)
	// expect lowercase + placeholder for mixed_Key
	if !findWarning(warnings, "mixed_Key", "lowercase") {
		t.Errorf("expected lowercase warning")
	}
	if !findWarning(warnings, "mixed_Key", "unresolved placeholder") {
		t.Errorf("expected placeholder warning")
	}
}

func TestLint_EmptyMap(t *testing.T) {
	warnings := linter.Lint(map[string]string{})
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for empty map, got %d", len(warnings))
	}
}
