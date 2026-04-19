package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// RollbackResult describes the outcome of a rollback operation.
type RollbackResult struct {
	RestoredChecksum string
	RestoredAt       time.Time
	PreviousChecksum string
}

// Rollback reverts the snapshot store to the entry immediately preceding
// the one identified by ref (checksum prefix or tag). The superseded entry
// is not deleted — only the "latest" pointer moves.
func Rollback(path, ref string) (*RollbackResult, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}
	if ref == "" {
		return nil, errors.New("ref is required")
	}

	entries, err := Load(path)
	if err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}
	if len(entries) == 0 {
		return nil, errors.New("snapshot is empty")
	}

	// Resolve ref to an index.
	targetIdx := -1
	for i, e := range entries {
		if e.Checksum == ref || (len(ref) >= 6 && len(e.Checksum) >= len(ref) && e.Checksum[:len(ref)] == ref) {
			targetIdx = i
			break
		}
		for _, t := range e.Tags {
			if t == ref {
				targetIdx = i
				break
			}
		}
		if targetIdx >= 0 {
			break
		}
	}
	if targetIdx < 0 {
		return nil, fmt.Errorf("ref %q not found in snapshot", ref)
	}
	if targetIdx == 0 {
		return nil, errors.New("no previous entry exists before the referenced snapshot")
	}

	current := entries[targetIdx]
	previous := entries[targetIdx-1]

	// Re-order: move previous entry to end so Latest() returns it.
	updated := make([]Entry, 0, len(entries))
	for i, e := range entries {
		if i != targetIdx-1 {
			updated = append(updated, e)
		}
	}
	updated = append(updated, previous)

	if err := Save(path, updated); err != nil {
		return nil, fmt.Errorf("save snapshot: %w", err)
	}

	return &RollbackResult{
		RestoredChecksum: previous.Checksum,
		RestoredAt:       time.Now().UTC(),
		PreviousChecksum: current.Checksum,
	}, nil
}
