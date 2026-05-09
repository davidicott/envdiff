package validator_test

import (
	"testing"

	"github.com/user/envdiff/internal/validator"
)

func TestValidate_ValidMap(t *testing.T) {
	env := map[string]string{
		"APP_ENV":       "production",
		"DATABASE_URL":  "postgres://localhost/db",
		"_PRIVATE_KEY":  "secret",
		"Port9000":      "9000",
	}
	issues := validator.Validate(env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	env := map[string]string{
		"": "some-value",
	}
	issues := validator.Validate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Message != "key must not be empty" {
		t.Errorf("unexpected message: %s", issues[0].Message)
	}
}

func TestValidate_InvalidKeyChars(t *testing.T) {
	env := map[string]string{
		"1INVALID":  "value",
		"BAD-KEY":   "value",
		"GOOD_KEY":  "value",
	}
	issues := validator.Validate(env)
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d: %v", len(issues), issues)
	}
	for _, iss := range issues {
		if iss.Key == "GOOD_KEY" {
			t.Errorf("GOOD_KEY should not produce an issue")
		}
	}
}

func TestValidate_NewlineInValue(t *testing.T) {
	env := map[string]string{
		"MULTILINE": "line1\nline2",
	}
	issues := validator.Validate(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "MULTILINE" {
		t.Errorf("expected issue on MULTILINE, got %s", issues[0].Key)
	}
}

func TestValidate_EmptyMap(t *testing.T) {
	issues := validator.Validate(map[string]string{})
	if len(issues) != 0 {
		t.Errorf("expected no issues for empty map, got %d", len(issues))
	}
}

func TestIssue_Error(t *testing.T) {
	issue := validator.Issue{Key: "BAD KEY", Value: "", Message: "key contains invalid characters"}
	got := issue.Error()
	want := `key "BAD KEY": key contains invalid characters`
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}
