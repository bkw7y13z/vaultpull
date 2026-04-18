package snapshot

import (
	"fmt"
	"os"
	"sort"
	"time"
)

// Stats holds aggregate statistics for a snapshot file.
type Stats struct {
	TotalEntries  int
	UniqueKeys    int
	TaggedEntries int
	PinnedEntries int
	OldestAt      time.Time
	NewestAt      time.Time
	TopKeys       []string // top 5 most frequently appearing keys
}

// ComputeStats loads a snapshot file and returns aggregate statistics.
func ComputeStats(path string) (*Stats, error) {
	if path == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}

	entries, err := Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("snapshot file not found: %s", path)
		}
		return nil, err
	}

	if len(entries) == 0 {
		return &Stats{}, nil
	}

	keyCount := map[string]int{}
	keySet := map[string]struct{}{}
	tagged := 0
	pinned := 0
	oldest := entries[0].CreatedAt
	newest := entries[0].CreatedAt

	for _, e := range entries {
		if e.Tag != "" {
			tagged++
		}
		if e.Pinned {
			pinned++
		}
		if e.CreatedAt.Before(oldest) {
			oldest = e.CreatedAt
		}
		if e.CreatedAt.After(newest) {
			newest = e.CreatedAt
		}
		for _, k := range e.Keys {
			keyCount[k]++
			keySet[k] = struct{}{}
		}
	}

	type kv struct {
		Key   string
		Count int
	}
	var sorted []kv
	for k, c := range keyCount {
		sorted = append(sorted, kv{k, c})
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Count != sorted[j].Count {
			return sorted[i].Count > sorted[j].Count
		}
		return sorted[i].Key < sorted[j].Key
	})

	top := []string{}
	for i := 0; i < len(sorted) && i < 5; i++ {
		top = append(top, sorted[i].Key)
	}

	return &Stats{
		TotalEntries:  len(entries),
		UniqueKeys:    len(keySet),
		TaggedEntries: tagged,
		PinnedEntries: pinned,
		OldestAt:      oldest,
		NewestAt:      newest,
		TopKeys:       top,
	}, nil
}
