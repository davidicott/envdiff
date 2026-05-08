// Package reporter formats and writes envdiff comparison results
// to an output stream.
//
// It supports rendering diff.Result slices into human-readable text,
// grouping entries by status: missing keys, extra keys, and mismatched values.
//
// The three result categories are:
//
//   - Missing: keys present in the reference file but absent in the target
//   - Extra: keys present in the target file but absent in the reference
//   - Mismatch: keys present in both files but with differing values
//
// Example usage:
//
//	results := diff.Compare(mapA, mapB)
//	reporter.Report(os.Stdout, results, ".env", ".env.production")
//
// Output is sorted alphabetically by key for consistent, readable reports.
package reporter
