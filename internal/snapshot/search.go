package snapshot

import (
	"errors"
	"strings"
	"time"
)

// SearchOptions defines criteria for filtering snapshot entries.
type SearchOptions struct {
	KeyContains   string
	TagEquals     string
	After         time.Time
	Before        time.Time
	PinnedOnly    bool
}

// SearchResult holds a matched entry and its index.
type SearchResult struct {
	Index int
	Entry Entry
}

// Search filters snapshot entries based on the provided options.
func Search(path string, opts SearchOptions) ([]SearchResult, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}

	snap, err := Load(path)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for i, e := range snap.Entries {
		if !opts.After.IsZero() && !e.CreatedAt.After(opts.After) {
			continue
		}
		if !opts.Before.IsZero() && !e.CreatedAt.Before(opts.Before) {
			continue
		}
		if opts.TagEquals != "" && e.Tag != opts.TagEquals {
			continue
		}
		if opts.PinnedOnly && !e.Pinned {
			continue
		}
		if opts.KeyContains != "" {
			matched := false
			for _, k := range e.Keys {
				if strings.Contains(k, opts.KeyContains) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
		}
		results = append(results, SearchResult{Index: i, Entry: e})
	}
	return results, nil
}
