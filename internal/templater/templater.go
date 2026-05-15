// Package templater provides functionality for rendering .env templates
// by substituting known values and flagging unresolved placeholders.
package templater

import (
	"fmt"
	"regexp"
	"strings"
)

// placeholderRe matches ${VAR_NAME} style placeholders.
var placeholderRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Result holds the output of rendering a single key's value.
type Result struct {
	Key        string
	Original   string
	Rendered   string
	Unresolved []string
}

// Render takes a template map (values may contain ${VAR} placeholders) and a
// values map used for substitution. It returns a Result for every key in the
// template map.
func Render(template map[string]string, values map[string]string) []Result {
	results := make([]Result, 0, len(template))

	for key, raw := range template {
		result := Result{
			Key:      key,
			Original: raw,
		}

		rendered, unresolved := substitute(raw, values)
		result.Rendered = rendered
		result.Unresolved = unresolved

		results = append(results, result)
	}

	return results
}

// substitute replaces all ${VAR} placeholders in s using the provided values
// map. It returns the substituted string and a slice of placeholder names that
// could not be resolved.
func substitute(s string, values map[string]string) (string, []string) {
	var unresolved []string

	out := placeholderRe.ReplaceAllStringFunc(s, func(match string) string {
		name := placeholderRe.FindStringSubmatch(match)[1]
		if val, ok := values[name]; ok {
			return val
		}
		unresolved = append(unresolved, name)
		return match // leave original placeholder intact
	})

	return out, unresolved
}

// ToMap converts a slice of Results into a plain key→rendered-value map,
// skipping entries that still have unresolved placeholders.
func ToMap(results []Result) (map[string]string, []string) {
	out := make(map[string]string, len(results))
	var warnings []string

	for _, r := range results {
		if len(r.Unresolved) > 0 {
			warnings = append(warnings,
				fmt.Sprintf("key %q has unresolved placeholders: %s",
					r.Key, strings.Join(r.Unresolved, ", ")))
			continue
		}
		out[r.Key] = r.Rendered
	}

	return out, warnings
}
