// Package exporter provides utilities for exporting envdiff results
// to external file formats.
//
// Supported formats:
//
//   - FormatJSON: newline-delimited JSON array with key, status, value_a, value_b fields.
//   - FormatCSV:  comma-separated values with a header row.
//
// Example usage:
//
//	f, _ := os.Create("results.csv")
//	defer f.Close()
//	exporter.Export(f, results, exporter.FormatCSV)
package exporter
