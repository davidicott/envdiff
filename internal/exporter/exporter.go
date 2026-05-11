// Package exporter provides functionality to export diff results
// to various file formats such as JSON and CSV.
package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the export format type.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Export writes the given diff results to w in the specified format.
// Returns an error if the format is unsupported or writing fails.
func Export(w io.Writer, results []diff.Result, format Format) error {
	switch format {
	case FormatJSON:
		return exportJSON(w, results)
	case FormatCSV:
		return exportCSV(w, results)
	default:
		return fmt.Errorf("unsupported export format: %q", format)
	}
}

func exportJSON(w io.Writer, results []diff.Result) error {
	type row struct {
		Key    string `json:"key"`
		Status string `json:"status"`
		ValueA string `json:"value_a,omitempty"`
		ValueB string `json:"value_b,omitempty"`
	}

	rows := make([]row, 0, len(results))
	for _, r := range results {
		rows = append(rows, row{
			Key:    r.Key,
			Status: string(r.Status),
			ValueA: r.ValueA,
			ValueB: r.ValueB,
		})
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

func exportCSV(w io.Writer, results []diff.Result) error {
	cw := csv.NewWriter(w)

	if err := cw.Write([]string{"key", "status", "value_a", "value_b"}); err != nil {
		return fmt.Errorf("writing csv header: %w", err)
	}

	for _, r := range results {
		record := []string{
			r.Key,
			strings.ToLower(string(r.Status)),
			r.ValueA,
			r.ValueB,
		}
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("writing csv record: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}
