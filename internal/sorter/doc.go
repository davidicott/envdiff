// Package sorter provides utilities for ordering diff results
// produced by the diff package.
//
// Results can be sorted by key name (alphabetically) or by status
// (missing-in-b, missing-in-a, mismatch, match), with a secondary
// alphabetical sort by key when statuses are equal.
//
// Example usage:
//
//	sorted := sorter.Results(results, sorter.SortByKey)
//	sortedByStatus := sorter.Results(results, sorter.SortByStatus)
package sorter
