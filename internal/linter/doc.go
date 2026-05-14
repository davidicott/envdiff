// Package linter performs heuristic analysis on parsed .env key-value maps.
//
// It is intentionally separate from validation (which enforces structural
// correctness) and diff (which compares two environments). The linter
// surfaces style and safety concerns within a single file, such as:
//
//   - Keys that do not follow the UPPER_SNAKE_CASE convention
//   - Values that contain unresolved shell-style placeholders (${VAR})
//   - Unusually large values that may indicate accidental data inclusion
//
// Usage:
//
//	warnings := linter.Lint(envMap)
//	for _, w := range warnings {
//		fmt.Println(w.Key, ":", w.Message)
//	}
package linter
