package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/reporter"
)

func TestReport_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	reporter.Report(&buf, []diff.Result{}, "a.env", "b.env")
	out := buf.String()
	if !strings.Contains(out, "No differences found") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestReport_MissingInB(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMissingInB, ValueA: "localhost"},
	}
	reporter.Report(&buf, results, "a.env", "b.env")
	out := buf.String()
	if !strings.Contains(out, "Missing in b.env") {
		t.Errorf("expected missing section, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key DB_HOST in output, got: %s", out)
	}
}

func TestReport_Mismatched(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "PORT", Status: diff.StatusMismatched, ValueA: "3000", ValueB: "4000"},
	}
	reporter.Report(&buf, results, "a.env", "b.env")
	out := buf.String()
	if !strings.Contains(out, "Mismatched values") {
		t.Errorf("expected mismatch section, got: %s", out)
	}
	if !strings.Contains(out, "3000") || !strings.Contains(out, "4000") {
		t.Errorf("expected both values in output, got: %s", out)
	}
}

func TestReport_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "Z_KEY", Status: diff.StatusMissingInB, ValueA: "z"},
		{Key: "A_KEY", Status: diff.StatusMissingInB, ValueA: "a"},
	}
	reporter.Report(&buf, results, "a.env", "b.env")
	out := buf.String()
	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Errorf("expected A_KEY before Z_KEY in sorted output")
	}
}

func TestReport_TotalCount(t *testing.T) {
	var buf bytes.Buffer
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMissingInB},
		{Key: "B", Status: diff.StatusMissingInA},
		{Key: "C", Status: diff.StatusMismatched, ValueA: "1", ValueB: "2"},
	}
	reporter.Report(&buf, results, "a.env", "b.env")
	out := buf.String()
	if !strings.Contains(out, "Total issues: 3") {
		t.Errorf("expected total count of 3, got: %s", out)
	}
}
