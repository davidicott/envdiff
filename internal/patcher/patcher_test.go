package patcher_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/patcher"
)

func makeResult(key, a, b string, status diff.Status) diff.Result {
	return diff.Result{Key: key, ValueA: a, ValueB: b, Status: status}
}

func TestPatch_OKPassthrough(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	results := []diff.Result{makeResult("FOO", "bar", "bar", diff.StatusOK)}

	out, err := patcher.Patch(base, results, patcher.PreferSource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
}

func TestPatch_MissingInA_AddsKey(t *testing.T) {
	base := map[string]string{}
	results := []diff.Result{makeResult("NEW", "", "value", diff.StatusMissingInA)}

	out, err := patcher.Patch(base, results, patcher.PreferSource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %q", out["NEW"])
	}
}

func TestPatch_MissingInB_KeepsKey(t *testing.T) {
	base := map[string]string{"GONE": "still-here"}
	results := []diff.Result{makeResult("GONE", "still-here", "", diff.StatusMissingInB)}

	out, err := patcher.Patch(base, results, patcher.PreferSource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GONE"] != "still-here" {
		t.Errorf("expected GONE=still-here, got %q", out["GONE"])
	}
}

func TestPatch_Mismatch_PreferSource(t *testing.T) {
	base := map[string]string{"KEY": "source-val"}
	results := []diff.Result{makeResult("KEY", "source-val", "target-val", diff.StatusMismatched)}

	out, err := patcher.Patch(base, results, patcher.PreferSource)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "source-val" {
		t.Errorf("expected KEY=source-val, got %q", out["KEY"])
	}
}

func TestPatch_Mismatch_PreferTarget(t *testing.T) {
	base := map[string]string{"KEY": "source-val"}
	results := []diff.Result{makeResult("KEY", "source-val", "target-val", diff.StatusMismatched)}

	out, err := patcher.Patch(base, results, patcher.PreferTarget)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "target-val" {
		t.Errorf("expected KEY=target-val, got %q", out["KEY"])
	}
}

func TestPatch_Mismatch_ErrorOnMismatch(t *testing.T) {
	base := map[string]string{"KEY": "a"}
	results := []diff.Result{makeResult("KEY", "a", "b", diff.StatusMismatched)}

	_, err := patcher.Patch(base, results, patcher.ErrorOnMismatch)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
