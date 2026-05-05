// Package reporter provides formatting and output of diff results.
package reporter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable diff report to the given writer.
func Report(w io.Writer, results []diff.Result, fileA, fileB string) {
	if len(results) == 0 {
		fmt.Fprintf(w, "✓ No differences found between %s and %s\n", fileA, fileB)
		return
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	missing := filterByStatus(results, diff.StatusMissingInB)
	extra := filterByStatus(results, diff.StatusMissingInA)
	mismatched := filterByStatus(results, diff.StatusMismatched)

	fmt.Fprintf(w, "Comparing %s → %s\n", fileA, fileB)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	if len(missing) > 0 {
		fmt.Fprintf(w, "\nMissing in %s (%d):\n", fileB, len(missing))
		for _, r := range missing {
			fmt.Fprintf(w, "  - %s\n", r.Key)
		}
	}

	if len(extra) > 0 {
		fmt.Fprintf(w, "\nExtra in %s (%d):\n", fileB, len(extra))
		for _, r := range extra {
			fmt.Fprintf(w, "  + %s\n", r.Key)
		}
	}

	if len(mismatched) > 0 {
		fmt.Fprintf(w, "\nMismatched values (%d):\n", len(mismatched))
		for _, r := range mismatched {
			fmt.Fprintf(w, "  ~ %s: %q → %q\n", r.Key, r.ValueA, r.ValueB)
		}
	}

	fmt.Fprintf(w, "\nTotal issues: %d\n", len(results))
}

func filterByStatus(results []diff.Result, status diff.Status) []diff.Result {
	var filtered []diff.Result
	for _, r := range results {
		if r.Status == status {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
