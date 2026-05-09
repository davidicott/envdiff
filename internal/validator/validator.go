// Package validator provides utilities for validating parsed .env key-value pairs.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyPattern matches POSIX-compliant env var names: letters, digits, underscores,
// not starting with a digit.
var validKeyPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key     string
	Value   string
	Message string
}

func (i Issue) Error() string {
	return fmt.Sprintf("key %q: %s", i.Key, i.Message)
}

// Validate checks all key-value pairs in the provided map and returns a slice
// of Issues describing any problems found. An empty slice means no issues.
func Validate(env map[string]string) []Issue {
	var issues []Issue

	for k, v := range env {
		if k == "" {
			issues = append(issues, Issue{Key: k, Value: v, Message: "key must not be empty"})
			continue
		}

		if !validKeyPattern.MatchString(k) {
			issues = append(issues, Issue{
				Key:     k,
				Value:   v,
				Message: "key contains invalid characters (must match [A-Za-z_][A-Za-z0-9_]*)",
			})
		}

		if strings.Contains(v, "\n") {
			issues = append(issues, Issue{
				Key:     k,
				Value:   v,
				Message: "value contains unescaped newline",
			})
		}
	}

	return issues
}
