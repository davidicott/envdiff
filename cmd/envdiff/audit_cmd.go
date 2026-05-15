package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/parser"
)

func runAudit(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	includeUnchanged := fs.Bool("unchanged", false, "include unchanged keys in output")
	jsonOutput := fs.Bool("json", false, "output as JSON")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return fmt.Errorf("audit requires two .env files: envdiff audit <before> <after>")
	}

	beforeFile, afterFile := remaining[0], remaining[1]

	before, err := parser.ParseFile(beforeFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", beforeFile, err)
	}

	after, err := parser.ParseFile(afterFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", afterFile, err)
	}

	report := auditor.Audit(before, after, *includeUnchanged)

	if *jsonOutput {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	if len(report.Entries) == 0 {
		fmt.Fprintln(out, "No changes detected.")
		return nil
	}

	fmt.Fprintf(out, "Audit: %s → %s\n", beforeFile, afterFile)
	fmt.Fprintf(out, "Generated: %s\n\n", report.Timestamp.Format("2006-01-02T15:04:05Z"))

	for _, e := range report.Entries {
		switch e.Kind {
		case auditor.Added:
			fmt.Fprintf(out, "[+] %s = %q\n", e.Key, e.NewValue)
		case auditor.Removed:
			fmt.Fprintf(out, "[-] %s (was %q)\n", e.Key, e.OldValue)
		case auditor.Modified:
			fmt.Fprintf(out, "[~] %s: %q → %q\n", e.Key, e.OldValue, e.NewValue)
		case auditor.Unchanged:
			fmt.Fprintf(out, "[=] %s = %q\n", e.Key, e.NewValue)
		}
	}

	fmt.Fprintf(out, "\nTotal: %d change(s)\n", len(report.Entries))
	return nil
}

func init() {
	_ = os.Stderr // ensure os is used
}
