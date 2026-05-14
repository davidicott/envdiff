package main

import (
	"strings"

	"github.com/user/envdiff/internal/redactor"
)

// buildRedactor constructs a Redactor from CLI flags.
// If noRedact is true, a Redactor that never masks anything is returned.
// If extraPatterns is non-empty, those patterns are merged with the defaults.
func buildRedactor(noRedact bool, extraPatterns string) *redactor.Redactor {
	if noRedact {
		return redactor.New([]string{}) // empty patterns → nothing is sensitive
	}

	patterns := append([]string{}, redactor.DefaultSensitivePatterns...)

	if extraPatterns != "" {
		for _, p := range strings.Split(extraPatterns, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				patterns = append(patterns, p)
			}
		}
	}

	return redactor.New(patterns)
}
