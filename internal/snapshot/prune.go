package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// PruneOptions controls how old snapshot entries are removed.
type PruneOptions struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// KeepLast retains at least this many entries regardless of age.
	KeepLast int
}

// PruneResult summarises what was removed.
type PruneResult struct {
	Removed int
	Retained int
}

// Prune removes old entries from the snapshot file at path according to opts.
// The file is rewritten in place. If the file does not exist, Prune is a no-op.
func Prune(path string, opts PruneOptions) (PruneResult, error) {
	if path == "" {
		return PruneResult{}, fmt.Errorf("snapshot path must not be empty")
	}
	if opts.KeepLast < 0 {
		return PruneResult{}, fmt.Errorf("KeepLast must be >= 0")
	}

	s, err := Load(path)
	if err != nil {
		return PruneResult{}, fmt.Errorf("loading snapshot: %w", err)
	}

	// Sort entries newest-first so we can easily keep the N most recent.
	sort.Slice(s.Entries, func(i, j int) bool {
		return s.Entries[i].Timestamp.After(s.Entries[j].Timestamp)
	})

	cutoff := time.Time{}
	if opts.MaxAge > 0 {
		cutoff = time.Now().UTC().Add(-opts.MaxAge)
	}

	var kept []Entry
	for i, e := range s.Entries {
		withinKeepLast := opts.KeepLast > 0 && i < opts.KeepLast
		withinAge := cutoff.IsZero() || e.Timestamp.After(cutoff)
		if withinKeepLast || withinAge {
			kept = append(kept, e)
		}
	}

	removed := len(s.Entries) - len(kept)
	s.Entries = kept

	if err := s.save(path); err != nil {
		return PruneResult{}, fmt.Errorf("saving pruned snapshot: %w", err)
	}

	return PruneResult{Removed: removed, Retained: len(kept)}, nil
}
