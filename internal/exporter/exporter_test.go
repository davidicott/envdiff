package exporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/exporter"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_NAME", Status: diff.StatusOK, ValueA: "myapp", ValueB: "myapp"},
		{Key: "DB_HOST", Status: diff.StatusMissingInB, ValueA: "localhost", ValueB: ""},
		{Key: "API_KEY", Status: diff.StatusMismatched, ValueA: "abc", ValueB: "xyz"},
	}
}

func TestExport_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, makeResults(), exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"key": "APP_NAME"`) {
		t.Errorf("expected APP_NAME in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"status": "missing_in_b"`) {
		t.Errorf("expected missing_in_b status in JSON output, got: %s", out)
	}
	if !strings.Contains(out, `"value_a": "abc"`) {
		t.Errorf("expected value_a=abc in JSON output, got: %s", out)
	}
}

func TestExport_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, makeResults(), exporter.FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 4 { // header + 3 data rows
		t.Fatalf("expected 4 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "key,status,value_a,value_b" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "APP_NAME,ok") {
		t.Errorf("unexpected first data row: %s", lines[1])
	}
}

func TestExport_CSVFormat_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, []diff.Result{}, exporter.FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := exporter.Export(&buf, makeResults(), exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported export format") {
		t.Errorf("unexpected error message: %v", err)
	}
}
