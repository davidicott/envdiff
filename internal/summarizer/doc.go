// Package summarizer provides aggregation utilities for diff results.
//
// It computes counts of OK, missing, extra, and mismatched keys from a
// []diff.Result slice, and exposes a HasIssues helper to quickly determine
// whether any actionable differences were found.
package summarizer
