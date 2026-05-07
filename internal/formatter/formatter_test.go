package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/formatter"
)

var sampleResults = []diff.Result{
	{Key: "DB_HOST", Status: diff.StatusMismatch, ValueA: "localhost", ValueB: "prod-db"},
	{Key: "API_KEY", Status: diff.StatusMissingInB},
	{Key: "SECRET", Status: diff.StatusMissingInA},
}

func TestWrite_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, sampleResults, formatter.FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[MISSING_IN_B] API_KEY") {
		t.Errorf("expected MISSING_IN_B for API_KEY, got:\n%s", out)
	}
	if !strings.Contains(out, "[MISSING_IN_A] SECRET") {
		t.Errorf("expected MISSING_IN_A for SECRET, got:\n%s", out)
	}
	if !strings.Contains(out, "[MISMATCH]     DB_HOST") {
		t.Errorf("expected MISMATCH for DB_HOST, got:\n%s", out)
	}
}

func TestWrite_TextFormat_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	formatter.Write(&buf, sampleResults, formatter.FormatText)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	// API_KEY < DB_HOST < SECRET alphabetically
	if !strings.Contains(lines[0], "API_KEY") {
		t.Errorf("first line should be API_KEY, got: %s", lines[0])
	}
}

func TestWrite_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	if err := formatter.Write(&buf, sampleResults, formatter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result formatter.JSONResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.Total != 3 {
		t.Errorf("expected total=3, got %d", result.Total)
	}
	if len(result.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(result.Entries))
	}
}

func TestWrite_JSONFormat_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	formatter.Write(&buf, []diff.Result{}, formatter.FormatJSON)
	var result formatter.JSONResult
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result.Total != 0 {
		t.Errorf("expected total=0, got %d", result.Total)
	}
}

func TestWrite_DefaultFormat(t *testing.T) {
	var buf bytes.Buffer
	// unknown format should fall back to text
	if err := formatter.Write(&buf, sampleResults, formatter.Format("unknown")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected non-empty text output for unknown format")
	}
}
