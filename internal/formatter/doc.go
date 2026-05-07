// Package formatter handles rendering diff results in multiple output formats.
//
// Supported formats:
//
//   - text (default): human-readable, line-oriented output suitable for
//     terminals and CI logs.
//
//   - json: structured JSON output for programmatic consumption or
//     integration with other tools.
//
// Usage:
//
//	var buf bytes.Buffer
//	err := formatter.Write(&buf, results, formatter.FormatJSON)
//
The Write function accepts any io.Writer, making it easy to direct output
to stdout, a file, or an in-memory buffer for testing.
package formatter
