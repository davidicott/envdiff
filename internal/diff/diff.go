// Package diff provides functionality for comparing two parsed .env maps
// and reporting missing or mismatched keys between them.
package diff

// Result holds the outcome of comparing two environment files.
type Result struct {
	// MissingInB contains keys present in A but absent in B.
	MissingInB []string
	// MissingInA contains keys present in B but absent in A.
	MissingInA []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey represents a key whose value differs between two env files.
type MismatchedKey struct {
	Key    string
	ValueA string
	ValueB string
}

// Compare takes two maps (parsed .env files) and returns a Result describing
// the differences between them. The maps are keyed by variable name.
func Compare(a, b map[string]string) Result {
	var result Result

	for key, valA := range a {
		valB, exists := b[key]
		if !exists {
			result.MissingInB = append(result.MissingInB, key)
			continue
		}
		if valA != valB {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:    key,
				ValueA: valA,
				ValueB: valB,
			})
		}
	}

	for key := range b {
		if _, exists := a[key]; !exists {
			result.MissingInA = append(result.MissingInA, key)
		}
	}

	return result
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInA) > 0 || len(r.MissingInB) > 0 || len(r.Mismatched) > 0
}
