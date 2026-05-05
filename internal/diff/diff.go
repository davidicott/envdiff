// Package diff compares two maps of environment variables and returns
// a list of differences including missing and mismatched keys.
package diff

// Status represents the type of difference found for a key.
type Status string

const (
	// StatusMissingInB indicates the key exists in A but not in B.
	StatusMissingInB Status = "missing_in_b"
	// StatusMissingInA indicates the key exists in B but not in A.
	StatusMissingInA Status = "missing_in_a"
	// StatusMismatched indicates the key exists in both but values differ.
	StatusMismatched Status = "mismatched"
)

// Result holds the comparison result for a single key.
type Result struct {
	Key     string
	Status  Status
	ValueA  string
	ValueB  string
}

// Compare compares two env maps and returns a slice of Result
// describing all differences between them.
func Compare(a, b map[string]string) []Result {
	var results []Result

	for key, valA := range a {
		valB, ok := b[key]
		if !ok {
			results = append(results, Result{
				Key:    key,
				Status: StatusMissingInB,
				ValueA: valA,
			})
			continue
		}
		if valA != valB {
			results = append(results, Result{
				Key:    key,
				Status: StatusMismatched,
				ValueA: valA,
				ValueB: valB,
			})
		}
	}

	for key, valB := range b {
		if _, ok := a[key]; !ok {
			results = append(results, Result{
				Key:    key,
				Status: StatusMissingInA,
				ValueB: valB,
			})
		}
	}

	return results
}
