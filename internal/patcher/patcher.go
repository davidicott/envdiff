// Package patcher provides functionality to apply a diff result set
// back to an env map, producing a patched output map.
package patcher

import (
	"fmt"

	"github.com/user/envdiff/internal/diff"
)

// Strategy controls how conflicts are resolved during patching.
type Strategy int

const (
	// PreferSource keeps the source (A) value on mismatch.
	PreferSource Strategy = iota
	// PreferTarget keeps the target (B) value on mismatch.
	PreferTarget
	// ErrorOnMismatch returns an error when a mismatch is encountered.
	ErrorOnMismatch
)

// Patch applies the given diff results to base, returning a new map that
// reconciles missing and mismatched keys according to the chosen strategy.
//
// - MissingInB entries are added to the result using the value from base.
// - MissingInA entries are added to the result using the value from results.
// - Mismatched entries are resolved according to strategy.
// - OK entries are passed through unchanged.
func Patch(base map[string]string, results []diff.Result, strategy Strategy) (map[string]string, error) {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}

	for _, r := range results {
		switch r.Status {
		case diff.StatusOK:
			out[r.Key] = r.ValueA

		case diff.StatusMissingInB:
			// Key exists in A but not B — already in out from base copy.
			out[r.Key] = r.ValueA

		case diff.StatusMissingInA:
			// Key exists in B but not A — add it using ValueB.
			out[r.Key] = r.ValueB

		case diff.StatusMismatched:
			switch strategy {
			case PreferSource:
				out[r.Key] = r.ValueA
			case PreferTarget:
				out[r.Key] = r.ValueB
			case ErrorOnMismatch:
				return nil, fmt.Errorf("patcher: mismatch on key %q (source=%q, target=%q)", r.Key, r.ValueA, r.ValueB)
			}
		}
	}

	return out, nil
}
