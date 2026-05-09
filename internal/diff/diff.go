// Package diff compares two maps of environment variables and returns
// a list of results describing keys that are missing or mismatched.
package diff

// Status constants describe the relationship of a key between two env maps.
const (
	StatusMatch      = "match"
	StatusMismatch   = "mismatch"
	StatusMissingInA = "missing_in_a"
	StatusMissingInB = "missing_in_b"
)

// Result holds the comparison outcome for a single environment key.
type Result struct {
	Key      string
	ValueA   string
	ValueB   string
	Status   string
}

// Compare takes two maps representing environment files A and B and
// returns a slice of Result describing every key found in either map.
func Compare(a, b map[string]string) []Result {
	seen := make(map[string]bool)
	var results []Result

	for k, va := range a {
		seen[k] = true
		if vb, ok := b[k]; !ok {
			results = append(results, Result{
				Key:    k,
				ValueA: va,
				Status: StatusMissingInB,
			})
		} else if va != vb {
			results = append(results, Result{
				Key:    k,
				ValueA: va,
				ValueB: vb,
				Status: StatusMismatch,
			})
		} else {
			results = append(results, Result{
				Key:    k,
				ValueA: va,
				ValueB: vb,
				Status: StatusMatch,
			})
		}
	}

	for k, vb := range b {
		if !seen[k] {
			results = append(results, Result{
				Key:    k,
				ValueB: vb,
				Status: StatusMissingInA,
			})
		}
	}

	return results
}
