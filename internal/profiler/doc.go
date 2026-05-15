// Package profiler analyses key distribution across multiple named
// environment maps. It reports how many environments each key appears
// in and surfaces "orphan" keys that exist in only one environment,
// which often indicate missing entries or copy-paste oversights.
//
// Example usage:
//
//	envs := map[string]map[string]string{
//		"production": prodMap,
//		"staging":    stagingMap,
//	}
//	profiles := profiler.Profile(envs)
//	orphanKeys := profiler.Orphans(profiles)
package profiler
