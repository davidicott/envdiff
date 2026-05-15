package auditor_test

import (
	"testing"

	"github.com/user/envdiff/internal/auditor"
)

func TestAudit_Added(t *testing.T) {
	before := map[string]string{"FOO": "bar"}
	after := map[string]string{"FOO": "bar", "BAZ": "qux"}

	report := auditor.Audit(before, after, false)

	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	if report.Entries[0].Kind != auditor.Added || report.Entries[0].Key != "BAZ" {
		t.Errorf("unexpected entry: %+v", report.Entries[0])
	}
}

func TestAudit_Removed(t *testing.T) {
	before := map[string]string{"FOO": "bar", "OLD": "val"}
	after := map[string]string{"FOO": "bar"}

	report := auditor.Audit(before, after, false)

	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	if report.Entries[0].Kind != auditor.Removed || report.Entries[0].Key != "OLD" {
		t.Errorf("unexpected entry: %+v", report.Entries[0])
	}
	if report.Entries[0].OldValue != "val" {
		t.Errorf("expected OldValue 'val', got %q", report.Entries[0].OldValue)
	}
}

func TestAudit_Modified(t *testing.T) {
	before := map[string]string{"KEY": "old"}
	after := map[string]string{"KEY": "new"}

	report := auditor.Audit(before, after, false)

	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	e := report.Entries[0]
	if e.Kind != auditor.Modified || e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestAudit_Unchanged_Excluded(t *testing.T) {
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "1", "B": "2"}

	report := auditor.Audit(before, after, false)

	if len(report.Entries) != 0 {
		t.Errorf("expected 0 entries when includeUnchanged=false, got %d", len(report.Entries))
	}
}

func TestAudit_Unchanged_Included(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "1"}

	report := auditor.Audit(before, after, true)

	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	if report.Entries[0].Kind != auditor.Unchanged {
		t.Errorf("expected Unchanged, got %s", report.Entries[0].Kind)
	}
}

func TestAudit_SortedByKey(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MANGO": "m"}

	report := auditor.Audit(before, after, false)

	keys := []string{report.Entries[0].Key, report.Entries[1].Key, report.Entries[2].Key}
	if keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("entries not sorted: %v", keys)
	}
}

func TestAudit_TimestampSet(t *testing.T) {
	report := auditor.Audit(map[string]string{}, map[string]string{}, false)
	if report.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}
