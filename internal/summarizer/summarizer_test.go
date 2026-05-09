package summarizer_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/summarizer"
)

func makeResults(statuses ...diff.Status) []diff.Result {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{Key: "KEY", Status: s}
	}
	return results
}

func TestSummarize_Empty(t *testing.T) {
	s := summarizer.Summarize(nil)
	if s.Total != 0 || s.HasIssues() {
		t.Errorf("expected empty summary, got %+v", s)
	}
}

func TestSummarize_AllOK(t *testing.T) {
	s := summarizer.Summarize(makeResults(diff.StatusOK, diff.StatusOK))
	if s.Total != 2 || s.OK != 2 || s.HasIssues() {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarize_MissingKeys(t *testing.T) {
	s := summarizer.Summarize(makeResults(diff.StatusMissingInB, diff.StatusMissingInB))
	if s.Missing != 2 || !s.HasIssues() {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarize_ExtraKeys(t *testing.T) {
	s := summarizer.Summarize(makeResults(diff.StatusMissingInA))
	if s.Extra != 1 || !s.HasIssues() {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestSummarize_Mixed(t *testing.T) {
	results := makeResults(
		diff.StatusOK,
		diff.StatusMissingInB,
		diff.StatusMissingInA,
		diff.StatusMismatch,
		diff.StatusMismatch,
	)
	s := summarizer.Summarize(results)
	if s.Total != 5 || s.OK != 1 || s.Missing != 1 || s.Extra != 1 || s.Mismatched != 2 {
		t.Errorf("unexpected summary: %+v", s)
	}
	if !s.HasIssues() {
		t.Error("expected HasIssues to be true")
	}
}

func TestHasIssues_FalseWhenAllOK(t *testing.T) {
	s := summarizer.Summary{Total: 3, OK: 3}
	if s.HasIssues() {
		t.Error("expected HasIssues to be false")
	}
}
