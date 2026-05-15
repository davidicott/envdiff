// Package profiler provides key usage profiling across multiple env maps,
// tracking which keys appear in which environments and how frequently.
package profiler

import "sort"

// KeyProfile holds statistics for a single key across environments.
type KeyProfile struct {
	Key          string
	Environments []string
	Count        int
}

// Profile returns a slice of KeyProfile entries describing how often and
// in which environments each key appears. Results are sorted by key name.
func Profile(envs map[string]map[string]string) []KeyProfile {
	type entry struct {
		envs map[string]bool
	}

	index := make(map[string]*entry)

	for envName, kv := range envs {
		for key := range kv {
			if _, ok := index[key]; !ok {
				index[key] = &entry{envs: make(map[string]bool)}
			}
			index[key].envs[envName] = true
		}
	}

	results := make([]KeyProfile, 0, len(index))
	for key, e := range index {
		envList := make([]string, 0, len(e.envs))
		for env := range e.envs {
			envList = append(envList, env)
		}
		sort.Strings(envList)
		results = append(results, KeyProfile{
			Key:          key,
			Environments: envList,
			Count:        len(envList),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Key < results[j].Key
	})

	return results
}

// Orphans returns keys that appear in only one environment out of all
// provided environments, indicating they may be forgotten or misplaced.
func Orphans(profiles []KeyProfile) []KeyProfile {
	var out []KeyProfile
	for _, p := range profiles {
		if p.Count == 1 {
			out = append(out, p)
		}
	}
	return out
}
