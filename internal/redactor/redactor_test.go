package redactor_test

import (
	"testing"

	"github.com/user/envdiff/internal/redactor"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	r := redactor.New(nil)

	sensitive := []string{"DB_PASSWORD", "API_KEY", "SECRET_TOKEN", "AUTH_TOKEN", "PRIVATE_KEY", "AWS_SECRET"}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}

	safe := []string{"APP_NAME", "PORT", "LOG_LEVEL", "DATABASE_HOST"}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestIsSensitive_CustomPatterns(t *testing.T) {
	r := redactor.New([]string{"INTERNAL", "HIDDEN"})

	if !r.IsSensitive("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to be sensitive")
	}
	if !r.IsSensitive("hidden_key") {
		t.Error("expected hidden_key to be sensitive (case-insensitive)")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom patterns")
	}
}

func TestRedact_MasksSensitiveValues(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"PORT":        "8080",
	}

	result := r.Redact(env)

	if result["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME should be unchanged, got %q", result["APP_NAME"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted, got %q", result["API_KEY"])
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	r := redactor.New(nil)
	env := map[string]string{
		"DB_PASSWORD": "original",
	}

	_ = r.Redact(env)

	if env["DB_PASSWORD"] != "original" {
		t.Error("Redact must not mutate the original map")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	r := redactor.New(nil)
	result := r.Redact(map[string]string{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
