// Package auditor compares two snapshots of environment variable maps
// and produces a structured audit report describing what changed.
//
// It detects four kinds of changes:
//
//   - Added: key present in after but not in before
//   - Removed: key present in before but not in after
//   - Modified: key present in both but with different values
//   - Unchanged: key present in both with identical values (opt-in)
//
// Usage:
//
//	report := auditor.Audit(before, after, includeUnchanged)
//	for _, entry := range report.Entries {
//		fmt.Printf("%s: %s\n", entry.Key, entry.Kind)
//	}
package auditor
