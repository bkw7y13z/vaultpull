package snapshot

import "sort"

// DiffResult holds the changes between two secret maps.
type DiffResult struct {
	Added   []string
	Removed []string
	Changed []string
}

// HasChanges returns true if there are any differences.
func (d DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Diff compares two secret maps (previous vs current) and returns a DiffResult
// describing which keys were added, removed, or changed.
func Diff(previous, current map[string]string) DiffResult {
	result := DiffResult{}

	// Find added and changed keys.
	for key, curVal := range current {
		prevVal, exists := previous[key]
		if !exists {
			result.Added = append(result.Added, key)
		} else if prevVal != curVal {
			result.Changed = append(result.Changed, key)
		}
	}

	// Find removed keys.
	for key := range previous {
		if _, exists := current[key]; !exists {
			result.Removed = append(result.Removed, key)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Changed)

	return result
}
