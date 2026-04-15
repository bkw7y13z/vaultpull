package snapshot

import (
	"fmt"
	"strings"
	"time"
)

// CompareResult holds the result of comparing two snapshot entries.
type CompareResult struct {
	FromChecksum string
	ToChecksum   string
	FromTime     time.Time
	ToTime       time.Time
	Diff         DiffResult
	Unchanged    bool
}

// String returns a human-readable summary of the comparison.
func (r CompareResult) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "From: %s (%s)\n", r.FromChecksum[:8], r.FromTime.Format(time.RFC3339))
	fmt.Fprintf(&sb, "To:   %s (%s)\n", r.ToChecksum[:8], r.ToTime.Format(time.RFC3339))
	if r.Unchanged {
		sb.WriteString("No changes detected.\n")
		return sb.String()
	}
	if len(r.Diff.Added) > 0 {
		fmt.Fprintf(&sb, "Added (%d):   %s\n", len(r.Diff.Added), strings.Join(r.Diff.Added, ", "))
	}
	if len(r.Diff.Removed) > 0 {
		fmt.Fprintf(&sb, "Removed (%d): %s\n", len(r.Diff.Removed), strings.Join(r.Diff.Removed, ", "))
	}
	if len(r.Diff.Changed) > 0 {
		fmt.Fprintf(&sb, "Changed (%d): %s\n", len(r.Diff.Changed), strings.Join(r.Diff.Changed, ", "))
	}
	return sb.String()
}

// Compare loads the snapshot at path and compares the two most recent entries
// identified by their checksums. If only one entry exists it returns an error.
func Compare(path, fromChecksum, toChecksum string) (CompareResult, error) {
	if path == "" {
		return CompareResult{}, fmt.Errorf("snapshot path must not be empty")
	}
	if fromChecksum == "" || toChecksum == "" {
		return CompareResult{}, fmt.Errorf("both from and to checksums must be provided")
	}

	snap, err := Load(path)
	if err != nil {
		return CompareResult{}, fmt.Errorf("loading snapshot: %w", err)
	}

	var from, to *Entry
	for i := range snap.Entries {
		e := &snap.Entries[i]
		if strings.HasPrefix(e.Checksum, fromChecksum) {
			from = e
		}
		if strings.HasPrefix(e.Checksum, toChecksum) {
			to = e
		}
	}
	if from == nil {
		return CompareResult{}, fmt.Errorf("entry with checksum prefix %q not found", fromChecksum)
	}
	if to == nil {
		return CompareResult{}, fmt.Errorf("entry with checksum prefix %q not found", toChecksum)
	}

	diff := Diff(from.Keys, to.Keys)
	unchanged := len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Changed) == 0

	return CompareResult{
		FromChecksum: from.Checksum,
		ToChecksum:   to.Checksum,
		FromTime:     from.Timestamp,
		ToTime:       to.Timestamp,
		Diff:         diff,
		Unchanged:    unchanged,
	}, nil
}
