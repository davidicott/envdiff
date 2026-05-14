package redactor

import "strings"

// DefaultSensitivePatterns is a list of key substrings that are considered sensitive.
var DefaultSensitivePatterns = []string{
	"PASSWORD",
	"SECRET",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"PASSWD",
}

const redactedValue = "[REDACTED]"

// Redactor masks sensitive values in env maps.
type Redactor struct {
	patterns []string
}

// New creates a Redactor with the given patterns.
// If no patterns are provided, DefaultSensitivePatterns are used.
func New(patterns []string) *Redactor {
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}
	upper := make([]string, len(patterns))
	for i, p := range patterns {
		upper[i] = strings.ToUpper(p)
	}
	return &Redactor{patterns: upper}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Redact returns a copy of the map with sensitive values replaced by [REDACTED].
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out
}
