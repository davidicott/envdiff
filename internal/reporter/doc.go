// Package reporter formats and writes envdiff comparison results
// to an output stream.
//
// It supports rendering diff.Result slices into human-readable text,
// grouping entries by status: missing keys, extra keys, and mismatched values.
//
// Example usage:
//
//	results := diff.Compare(mapA, mapB)
//	reporter.Report(os.Stdout, results, ".env", ".env.production")
//
// Output is sorted alphabetically by key for consistent, readable reports.
package reporter
