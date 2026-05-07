// Package formatter provides output formatting for envdiff results.
package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format type.
type Format string

const (
	// FormatText is the default human-readable text format.
	FormatText Format = "text"
	// FormatJSON outputs results as JSON.
	FormatJSON Format = "json"
)

// JSONResult is the structure used for JSON output.
type JSONResult struct {
	Total    int         `json:"total"`
	Entries  []JSONEntry `json:"entries"`
}

// JSONEntry represents a single diff entry in JSON output.
type JSONEntry struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	ValueA string `json:"value_a,omitempty"`
	ValueB string `json:"value_b,omitempty"`
}

// Write formats and writes diff results to w using the given format.
func Write(w io.Writer, results []diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return writeJSON(w, results)
	default:
		return writeText(w, results)
	}
}

func writeText(w io.Writer, results []diff.Result) error {
	sorted := make([]diff.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})
	for _, r := range sorted {
		switch r.Status {
		case diff.StatusMissingInB:
			fmt.Fprintf(w, "[MISSING_IN_B] %s\n", r.Key)
		case diff.StatusMissingInA:
			fmt.Fprintf(w, "[MISSING_IN_A] %s\n", r.Key)
		case diff.StatusMismatch:
			fmt.Fprintf(w, "[MISMATCH]     %s: %q != %q\n", r.Key, r.ValueA, r.ValueB)
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []diff.Result) error {
	entries := make([]JSONEntry, 0, len(results))
	for _, r := range results {
		entries = append(entries, JSONEntry{
			Key:    r.Key,
			Status: string(r.Status),
			ValueA: r.ValueA,
			ValueB: r.ValueB,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	out := JSONResult{Total: len(entries), Entries: entries}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
