package snapshot

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"
)

// Summary holds aggregated statistics about a snapshot history file.
type Summary struct {
	TotalEntries  int
	LatestChecksum string
	LatestAt      time.Time
	OldestAt      time.Time
	UniqueKeys    []string
}

// Summarize reads the snapshot file at path and returns a Summary.
func Summarize(path string) (*Summary, error) {
	if path == "" {
		return nil, fmt.Errorf("snapshot path must not be empty")
	}

	entries, err := Load(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Summary{}, nil
		}
		return nil, fmt.Errorf("loading snapshot: %w", err)
	}

	if len(entries) == 0 {
		return &Summary{}, nil
	}

	keySet := make(map[string]struct{})
	oldest := entries[0].CreatedAt
	newest := entries[0].CreatedAt

	for _, e := range entries {
		for _, k := range e.Keys {
			keySet[k] = struct{}{}
		}
		if e.CreatedAt.Before(oldest) {
			oldest = e.CreatedAt
		}
		if e.CreatedAt.After(newest) {
			newest = e.CreatedAt
		}
	}

	latest := entries[len(entries)-1]

	s := &Summary{
		TotalEntries:   len(entries),
		LatestChecksum: latest.Checksum,
		LatestAt:       newest,
		OldestAt:       oldest,
		UniqueKeys:     sortedKeys(keySet),
	}
	return s, nil
}

// Print writes a human-readable summary table to w.
func (s *Summary) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Total Entries:\t%d\n", s.TotalEntries)
	if s.TotalEntries > 0 {
		fmt.Fprintf(tw, "Latest Checksum:\t%s\n", s.LatestChecksum)
		fmt.Fprintf(tw, "Latest Snapshot:\t%s\n", s.LatestAt.Format(time.RFC3339))
		fmt.Fprintf(tw, "Oldest Snapshot:\t%s\n", s.OldestAt.Format(time.RFC3339))
		fmt.Fprintf(tw, "Unique Keys Seen:\t%d\n", len(s.UniqueKeys))
	}
	tw.Flush()
}
