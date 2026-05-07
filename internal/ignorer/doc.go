// Package ignorer provides an Ignorer type that manages a set of environment
// variable keys to be excluded when comparing .env files.
//
// Keys can be loaded from a plain-text ignore file (one key per line,
// lines starting with '#' are treated as comments) or added programmatically
// via Add.
//
// Example ignore file (.envignore):
//
//	# Secrets managed externally
//	AWS_SECRET_ACCESS_KEY
//	DATABASE_PASSWORD
//
// Usage:
//
//	ig := ignorer.New()
//	_ = ig.LoadFile(".envignore")
//	filtered := ig.Filter(parsedEnvMap)
package ignorer
