// Package auditor provides functionality for auditing .env file changes
// between two versions, producing a structured changelog of additions,
// removals, and modifications.
package auditor

import (
	"sort"
	"time"
)

// ChangeKind describes the type of change detected.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Entry represents a single audited change for a key.
type Entry struct {
	Key      string     `json:"key"`
	Kind     ChangeKind `json:"kind"`
	OldValue string     `json:"old_value,omitempty"`
	NewValue string     `json:"new_value,omitempty"`
}

// Report holds the full audit result.
type Report struct {
	Timestamp time.Time `json:"timestamp"`
	Entries   []Entry   `json:"entries"`
}

// Audit compares two env maps (before and after) and returns a Report
// describing all additions, removals, modifications, and unchanged keys.
func Audit(before, after map[string]string, includeUnchanged bool) Report {
	seen := make(map[string]bool)
	var entries []Entry

	for key, newVal := range after {
		seen[key] = true
		if oldVal, exists := before[key]; !exists {
			entries = append(entries, Entry{Key: key, Kind: Added, NewValue: newVal})
		} else if oldVal != newVal {
			entries = append(entries, Entry{Key: key, Kind: Modified, OldValue: oldVal, NewValue: newVal})
		} else if includeUnchanged {
			entries = append(entries, Entry{Key: key, Kind: Unchanged, OldValue: oldVal, NewValue: newVal})
		}
	}

	for key, oldVal := range before {
		if !seen[key] {
			entries = append(entries, Entry{Key: key, Kind: Removed, OldValue: oldVal})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return Report{
		Timestamp: time.Now().UTC(),
		Entries:   entries,
	}
}
