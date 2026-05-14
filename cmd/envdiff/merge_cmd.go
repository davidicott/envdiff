package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/envdiff/internal/merger"
	"github.com/user/envdiff/internal/parser"
)

// runMerge handles the `envdiff merge` sub-command.
// Usage: envdiff merge [--strategy=a|b|error] <fileA> <fileB>
func runMerge(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	strategyFlag := fs.String("strategy", "a", "conflict resolution strategy: a, b, or error")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 2 {
		return fmt.Errorf("merge requires two file arguments")
	}

	fileA, fileB := fs.Arg(0), fs.Arg(1)

	mapA, err := parser.ParseFile(fileA)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fileA, err)
	}
	mapB, err := parser.ParseFile(fileB)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fileB, err)
	}

	var strategy merger.Strategy
	switch *strategyFlag {
	case "a":
		strategy = merger.PreferA
	case "b":
		strategy = merger.PreferB
	case "error":
		strategy = merger.ErrorOnConflict
	default:
		return fmt.Errorf("unknown strategy %q: use a, b, or error", *strategyFlag)
	}

	res, err := merger.Merge(mapA, mapB, strategy)
	if err != nil {
		return err
	}

	if len(res.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "conflicts resolved (%d): %v\n", len(res.Conflicts), res.Conflicts)
	}

	keys := make([]string, 0, len(res.Merged))
	for k := range res.Merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(stdout, "%s=%s\n", k, res.Merged[k])
	}
	return nil
}
