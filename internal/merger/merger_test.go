package merger_test

import (
	"testing"

	"github.com/user/envdiff/internal/merger"
)

func TestMerge_DisjointMaps(t *testing.T) {
	a := map[string]string{"FOO": "1"}
	b := map[string]string{"BAR": "2"}
	res, err := merger.Merge(a, b, merger.PreferA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["FOO"] != "1" || res.Merged["BAR"] != "2" {
		t.Errorf("unexpected merged map: %v", res.Merged)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_PreferA(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	res, err := merger.Merge(a, b, merger.PreferA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "from-a" {
		t.Errorf("expected from-a, got %q", res.Merged["KEY"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "KEY" {
		t.Errorf("expected conflict on KEY, got %v", res.Conflicts)
	}
}

func TestMerge_PreferB(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}
	res, err := merger.Merge(a, b, merger.PreferB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Merged["KEY"] != "from-b" {
		t.Errorf("expected from-b, got %q", res.Merged["KEY"])
	}
}

func TestMerge_ErrorOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "val1"}
	b := map[string]string{"KEY": "val2"}
	_, err := merger.Merge(a, b, merger.ErrorOnConflict)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_SameValues_NoConflict(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}
	res, err := merger.Merge(a, b, merger.ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts for identical values")
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	res, err := merger.Merge(map[string]string{}, map[string]string{}, merger.PreferA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Merged) != 0 {
		t.Errorf("expected empty merged map")
	}
}
