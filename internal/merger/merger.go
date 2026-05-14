package merger

import (
	"fmt"
	"maps"
	"sort"
)

// Strategy defines how conflicting keys are resolved during a merge.
type Strategy int

const (
	// PreferA keeps the value from map A when a key exists in both.
	PreferA Strategy = iota
	// PreferB keeps the value from map B when a key exists in both.
	PreferB
	// ErrorOnConflict returns an error if a key exists in both maps with different values.
	ErrorOnConflict
)

// Result holds the merged environment map and metadata about the merge.
type Result struct {
	Merged    map[string]string
	Conflicts []string
}

// Merge combines two env maps according to the given strategy.
// Keys unique to either map are always included in the result.
func Merge(a, b map[string]string, strategy Strategy) (Result, error) {
	merged := make(map[string]string)
	maps.Copy(merged, a)

	var conflicts []string

	for k, vb := range b {
		va, exists := merged[k]
		if !exists {
			merged[k] = vb
			continue
		}
		if va == vb {
			continue
		}
		conflicts = append(conflicts, k)
		switch strategy {
		case PreferA:
			// keep va — already in merged
		case PreferB:
			merged[k] = vb
		case ErrorOnConflict:
			sort.Strings(conflicts)
			return Result{}, fmt.Errorf("merger: conflict on key %q: %q vs %q", k, va, vb)
		}
	}

	sort.Strings(conflicts)
	return Result{Merged: merged, Conflicts: conflicts}, nil
}
