package templater_test

import (
	"testing"

	"github.com/user/envdiff/internal/templater"
)

func findResult(results []templater.Result, key string) (templater.Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return templater.Result{}, false
}

func TestRender_NoPlaceholders(t *testing.T) {
	tmpl := map[string]string{"HOST": "localhost", "PORT": "5432"}
	values := map[string]string{}

	results := templater.Render(tmpl, values)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	r, ok := findResult(results, "HOST")
	if !ok || r.Rendered != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", r.Rendered)
	}
}

func TestRender_ResolvesPlaceholders(t *testing.T) {
	tmpl := map[string]string{"DSN": "postgres://${DB_USER}:${DB_PASS}@localhost/mydb"}
	values := map[string]string{"DB_USER": "admin", "DB_PASS": "secret"}

	results := templater.Render(tmpl, values)
	r, ok := findResult(results, "DSN")
	if !ok {
		t.Fatal("DSN result not found")
	}
	want := "postgres://admin:secret@localhost/mydb"
	if r.Rendered != want {
		t.Errorf("expected %q, got %q", want, r.Rendered)
	}
	if len(r.Unresolved) != 0 {
		t.Errorf("expected no unresolved, got %v", r.Unresolved)
	}
}

func TestRender_UnresolvedPlaceholder(t *testing.T) {
	tmpl := map[string]string{"URL": "https://${API_HOST}/path"}
	values := map[string]string{}

	results := templater.Render(tmpl, values)
	r, ok := findResult(results, "URL")
	if !ok {
		t.Fatal("URL result not found")
	}
	if r.Rendered != "https://${API_HOST}/path" {
		t.Errorf("expected original placeholder preserved, got %q", r.Rendered)
	}
	if len(r.Unresolved) != 1 || r.Unresolved[0] != "API_HOST" {
		t.Errorf("expected [API_HOST] unresolved, got %v", r.Unresolved)
	}
}

func TestRender_EmptyTemplate(t *testing.T) {
	results := templater.Render(map[string]string{}, map[string]string{"X": "1"})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestToMap_SkipsUnresolved(t *testing.T) {
	results := []templater.Result{
		{Key: "GOOD", Rendered: "value", Unresolved: nil},
		{Key: "BAD", Rendered: "${MISSING}", Unresolved: []string{"MISSING"}},
	}

	out, warnings := templater.ToMap(results)
	if _, ok := out["GOOD"]; !ok {
		t.Error("expected GOOD in output map")
	}
	if _, ok := out["BAD"]; ok {
		t.Error("BAD should be excluded from output map")
	}
	if len(warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(warnings))
	}
}
