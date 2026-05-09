package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/sorter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "ZEBRA", Status: diff.StatusMatch},
		{Key: "ALPHA", Status: diff.StatusMismatch},
		{Key: "MONGO", Status: diff.StatusMissingInB},
		{Key: "BETA", Status: diff.StatusMissingInA},
		{Key: "APPLE", Status: diff.StatusMatch},
	}
}

func TestResults_SortByKey(t *testing.T) {
	in := makeResults()
	out := sorter.Results(in, sorter.SortByKey)

	expected := []string{"ALPHA", "APPLE", "BETA", "MONGO", "ZEBRA"}
	for i, r := range out {
		if r.Key != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, r.Key, expected[i])
		}
	}
}

func TestResults_SortByStatus(t *testing.T) {
	in := makeResults()
	out := sorter.Results(in, sorter.SortByStatus)

	// First result should be a MissingInB, last should be a Match
	if out[0].Status != diff.StatusMissingInB {
		t.Errorf("expected first status %q, got %q", diff.StatusMissingInB, out[0].Status)
	}
	if out[len(out)-1].Status != diff.StatusMatch {
		t.Errorf("expected last status %q, got %q", diff.StatusMatch, out[len(out)-1].Status)
	}
}

func TestResults_SortByStatus_TiebreakByKey(t *testing.T) {
	in := []diff.Result{
		{Key: "Z_KEY", Status: diff.StatusMatch},
		{Key: "A_KEY", Status: diff.StatusMatch},
	}
	out := sorter.Results(in, sorter.SortByStatus)

	if out[0].Key != "A_KEY" {
		t.Errorf("expected A_KEY first, got %q", out[0].Key)
	}
}

func TestResults_DefaultFallback(t *testing.T) {
	in := makeResults()
	out := sorter.Results(in, sorter.SortBy("unknown"))

	// Should fall back to key sort
	if out[0].Key != "ALPHA" {
		t.Errorf("expected ALPHA first, got %q", out[0].Key)
	}
}

func TestResults_DoesNotMutateInput(t *testing.T) {
	in := makeResults()
	originalFirst := in[0].Key
	_ = sorter.Results(in, sorter.SortByKey)

	if in[0].Key != originalFirst {
		t.Error("input slice was mutated")
	}
}
