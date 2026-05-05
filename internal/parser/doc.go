// Package parser provides utilities for reading and parsing .env files
// into structured key-value maps.
//
// Supported .env format:
//
//	# This is a comment
//	KEY=value
//	KEY="quoted value"
//	KEY='single quoted'
//
// Lines beginning with '#' are treated as comments and skipped.
// Empty lines are ignored. Values may optionally be wrapped in single
// or double quotes, which will be stripped during parsing.
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package parser
