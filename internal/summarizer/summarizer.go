package summarizer

import (
	"github.com/user/envdiff/internal/diff"
)

// Summary holds aggregated counts of a diff result set.
type Summary struct {
	Total     int
	Missing   int
	Extra     int
	Mismatched int
	OK        int
}

// Summarize computes a Summary from a slice of diff.Result values.
func Summarize(results []diff.Result) Summary {
	s := Summary{
		Total: len(results),
	}
	for _, r := range results {
		switch r.Status {
		case diff.StatusMissingInB:
			s.Missing++
		case diff.StatusMissingInA:
			s.Extra++
		case diff.StatusMismatch:
			s.Mismatched++
		case diff.StatusOK:
			s.OK++
		}
	}
	return s
}

// HasIssues returns true when any keys are missing, extra, or mismatched.
func (s Summary) HasIssues() bool {
	return s.Missing > 0 || s.Extra > 0 || s.Mismatched > 0
}
