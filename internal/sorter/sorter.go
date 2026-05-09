// Package sorter provides utilities for sorting diff results
// by various criteria such as key name, status, or file order.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// SortBy defines the field to sort results by.
type SortBy string

const (
	// SortByKey sorts results alphabetically by key name.
	SortByKey SortBy = "key"
	// SortByStatus sorts results by their diff status.
	SortByStatus SortBy = "status"
)

// Results sorts a slice of diff.Result values by the given field.
// Falls back to SortByKey if an unknown field is provided.
func Results(results []diff.Result, by SortBy) []diff.Result {
	out := make([]diff.Result, len(results))
	copy(out, results)

	switch by {
	case SortByStatus:
		sort.SliceStable(out, func(i, j int) bool {
			if out[i].Status == out[j].Status {
				return out[i].Key < out[j].Key
			}
			return statusOrder(out[i].Status) < statusOrder(out[j].Status)
		})
	default:
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].Key < out[j].Key
		})
	}

	return out
}

// statusOrder returns a numeric priority for a given status string
// so that results can be sorted in a meaningful order.
func statusOrder(status string) int {
	switch status {
	case diff.StatusMissingInB:
		return 0
	case diff.StatusMissingInA:
		return 1
	case diff.StatusMismatch:
		return 2
	case diff.StatusMatch:
		return 3
	default:
		return 99
	}
}
