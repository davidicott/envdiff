// Package linter provides heuristic checks on parsed .env maps,
// flagging potentially problematic key-value pairs such as keys with
// lowercase characters, values that look like unresolved placeholders,
// or suspiciously long values.
package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Warning represents a single lint finding for a specific key.
type Warning struct {
	Key     string
	Message string
}

var (
	lowercaseKey    = regexp.MustCompile(`[a-z]`)
	placeholderVal  = regexp.MustCompile(`^\$\{.+\}$`)
)

const maxValueLen = 1024

// Lint inspects the provided env map and returns a slice of Warnings.
// It checks for:
//   - keys containing lowercase letters (convention violation)
//   - values that appear to be unresolved shell placeholders (${VAR})
//   - values exceeding maxValueLen characters
func Lint(env map[string]string) []Warning {
	var warnings []Warning

	for k, v := range env {
		if lowercaseKey.MatchString(k) {
			warnings = append(warnings, Warning{
				Key:     k,
				Message: fmt.Sprintf("key %q contains lowercase letters; prefer UPPER_SNAKE_CASE", k),
			})
		}

		if placeholderVal.MatchString(strings.TrimSpace(v)) {
			warnings = append(warnings, Warning{
				Key:     k,
				Message: fmt.Sprintf("value for %q looks like an unresolved placeholder: %q", k, v),
			})
		}

		if len(v) > maxValueLen {
			warnings = append(warnings, Warning{
				Key:     k,
				Message: fmt.Sprintf("value for %q exceeds %d characters (%d)", k, maxValueLen, len(v)),
			})
		}
	}

	return warnings
}
