// Package ignorer provides functionality to load and apply ignore rules
// for specific keys when comparing .env files.
package ignorer

import (
	"bufio"
	"os"
	"strings"
)

// Ignorer holds a set of key patterns to be excluded from diff results.
type Ignorer struct {
	keys map[string]struct{}
}

// New creates a new Ignorer with no ignored keys.
func New() *Ignorer {
	return &Ignorer{keys: make(map[string]struct{})}
}

// LoadFile reads an ignore file where each non-blank, non-comment line
// is treated as a key to ignore during comparison.
func (ig *Ignorer) LoadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ig.keys[line] = struct{}{}
	}
	return scanner.Err()
}

// Add adds a single key to the ignore set.
func (ig *Ignorer) Add(key string) {
	ig.keys[strings.TrimSpace(key)] = struct{}{}
}

// ShouldIgnore reports whether the given key should be excluded.
func (ig *Ignorer) ShouldIgnore(key string) bool {
	_, ok := ig.keys[key]
	return ok
}

// Filter removes ignored keys from the provided map, returning a new map.
func (ig *Ignorer) Filter(env map[string]string) map[string]string {
	result := make(map[string]string, len(env))
	for k, v := range env {
		if !ig.ShouldIgnore(k) {
			result[k] = v
		}
	}
	return result
}

// Keys returns a slice of all currently ignored keys.
func (ig *Ignorer) Keys() []string {
	out := make([]string, 0, len(ig.keys))
	for k := range ig.keys {
		out = append(out, k)
	}
	return out
}
