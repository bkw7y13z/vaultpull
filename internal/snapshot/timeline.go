package snapshot

import (
	"errors"
	"sort"
	"time"
)

// TimelineEntry represents a snapshot entry on a timeline view.
type TimelineEntry struct {
	Checksum  string
	Tag       string
	Note      string
	KeyCount  int
	CreatedAt time.Time
}

// TimelineOptions filters the timeline output.
type TimelineOptions struct {
	Since  time.Time
	Until  time.Time
	Tagged bool // only include entries with a tag
}

// Timeline returns snapshot entries ordered by time, optionally filtered.
func Timeline(path string, opts TimelineOptions) ([]TimelineEntry, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}

	store, err := Load(path)
	if err != nil {
		return nil, err
	}

	var entries []TimelineEntry
	for _, e := range store.Entries {
		if !opts.Since.IsZero() && e.CreatedAt.Before(opts.Since) {
			continue
		}
		if !opts.Until.IsZero() && e.CreatedAt.After(opts.Until) {
			continue
		}
		if opts.Tagged && e.Tag == "" {
			continue
		}
		entries = append(entries, TimelineEntry{
			Checksum:  e.Checksum,
			Tag:       e.Tag,
			Note:      e.Note,
			KeyCount:  len(e.Keys),
			CreatedAt: e.CreatedAt,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.Before(entries[j].CreatedAt)
	})

	return entries, nil
}
