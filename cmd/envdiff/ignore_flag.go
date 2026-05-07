package main

import (
	"fmt"
	"os"

	"github.com/yourusername/envdiff/internal/ignorer"
)

// loadIgnorer builds an Ignorer from the optional ignore file path and any
// inline keys supplied via --ignore-key flags.
// If ignoreFile is empty and ignoreKeys is nil, a no-op Ignorer is returned.
func loadIgnorer(ignoreFile string, ignoreKeys []string) (*ignorer.Ignorer, error) {
	ig := ignorer.New()

	if ignoreFile != "" {
		if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("ignore file not found: %s", ignoreFile)
		}
		if err := ig.LoadFile(ignoreFile); err != nil {
			return nil, fmt.Errorf("loading ignore file: %w", err)
		}
	}

	for _, key := range ignoreKeys {
		if key != "" {
			ig.Add(key)
		}
	}

	return ig, nil
}
