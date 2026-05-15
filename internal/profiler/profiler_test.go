package profiler_test

import (
	"testing"

	"github.com/user/envdiff/internal/profiler"
)

func TestProfile_SingleEnv(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"DB_HOST": "localhost", "PORT": "5432"},
	}
	results := profiler.Profile(envs)
	if len(results) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST first, got %s", results[0].Key)
	}
	if results[0].Count != 1 {
		t.Errorf("expected count 1, got %d", results[0].Count)
	}
}

func TestProfile_MultipleEnvs_SharedKey(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"DB_HOST": "prod-db", "SECRET": "abc"},
		"staging": {"DB_HOST": "staging-db"},
	}
	results := profiler.Profile(envs)

	var dbProfile *profiler.KeyProfile
	for i := range results {
		if results[i].Key == "DB_HOST" {
			dbProfile = &results[i]
		}
	}
	if dbProfile == nil {
		t.Fatal("DB_HOST profile not found")
	}
	if dbProfile.Count != 2 {
		t.Errorf("expected DB_HOST in 2 envs, got %d", dbProfile.Count)
	}
}

func TestProfile_SortedByKey(t *testing.T) {
	envs := map[string]map[string]string{
		"prod": {"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3"},
	}
	results := profiler.Profile(envs)
	keys := []string{results[0].Key, results[1].Key, results[2].Key}
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %s, got %s", i, expected[i], k)
		}
	}
}

func TestProfile_EmptyEnvs(t *testing.T) {
	results := profiler.Profile(map[string]map[string]string{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestOrphans_ReturnsOnlyCount1(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"SHARED": "a", "ONLY_PROD": "b"},
		"staging": {"SHARED": "a"},
	}
	profiles := profiler.Profile(envs)
	orphanList := profiler.Orphans(profiles)
	if len(orphanList) != 1 {
		t.Fatalf("expected 1 orphan, got %d", len(orphanList))
	}
	if orphanList[0].Key != "ONLY_PROD" {
		t.Errorf("expected ONLY_PROD orphan, got %s", orphanList[0].Key)
	}
}

func TestOrphans_NoneWhenAllShared(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"KEY": "1"},
		"staging": {"KEY": "2"},
	}
	profiles := profiler.Profile(envs)
	orphanList := profiler.Orphans(profiles)
	if len(orphanList) != 0 {
		t.Errorf("expected no orphans, got %d", len(orphanList))
	}
}
