// Command envdiff compares .env files across environments and flags
// missing or mismatched keys.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/reporter"
)

const usage = `envdiff — compare .env files across environments

Usage:
  envdiff [flags] <file-a> <file-b>

Flags:
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		showValues  = flag.Bool("values", false, "show actual values in output")
		statusFilter = flag.String("status", "", "filter by status: missing_in_b, missing_in_a, mismatched")
	)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		return fmt.Errorf("exactly two .env files required, got %d", len(args))
	}

	fileA, fileB := args[0], args[1]

	envA, err := parser.ParseFile(fileA)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", fileA, err)
	}

	envB, err := parser.ParseFile(fileB)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", fileB, err)
	}

	results := diff.Compare(envA, envB)

	opts := reporter.Options{
		ShowValues:   *showValues,
		StatusFilter: *statusFilter,
		LabelA:       fileA,
		LabelB:       fileB,
	}

	reporter.Report(os.Stdout, results, opts)

	for _, r := range results {
		if r.Status != diff.StatusMatch {
			os.Exit(1)
		}
	}

	return nil
}
