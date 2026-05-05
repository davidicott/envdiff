package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := diff.Compare(a, b)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "value"}
	b := map[string]string{"FOO": "bar"}

	result := diff.Compare(a, b)

	if len(result.MissingInB) != 1 || result.MissingInB[0] != "ONLY_A" {
		t.Errorf("expected MissingInB=[ONLY_A], got %v", result.MissingInB)
	}
	if len(result.MissingInA) != 0 {
		t.Errorf("expected no MissingInA, got %v", result.MissingInA)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "value"}

	result := diff.Compare(a, b)

	if len(result.MissingInA) != 1 || result.MissingInA[0] != "ONLY_B" {
		t.Errorf("expected MissingInA=[ONLY_B], got %v", result.MissingInA)
	}
	if len(result.MissingInB) != 0 {
		t.Errorf("expected no MissingInB, got %v", result.MissingInB)
	}
}

func TestCompare_Mismatched(t *testing.T) {
	a := map[string]string{"FOO": "original", "BAR": "same"}
	b := map[string]string{"FOO": "changed", "BAR": "same"}

	result := diff.Compare(a, b)

	if len(result.Mismatched) != 1 {
		t.Fatalf("expected 1 mismatch, got %d", len(result.Mismatched))
	}
	m := result.Mismatched[0]
	if m.Key != "FOO" || m.ValueA != "original" || m.ValueB != "changed" {
		t.Errorf("unexpected mismatch entry: %+v", m)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := diff.Compare(map[string]string{}, map[string]string{})
	if result.HasDifferences() {
		t.Errorf("expected no differences for empty maps")
	}
}

func TestHasDifferences(t *testing.T) {
	result := diff.Result{
		MissingInB: []string{"KEY"},
	}
	if !result.HasDifferences() {
		t.Error("expected HasDifferences to return true")
	}
}
