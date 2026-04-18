package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourceChecksum string
	NewChecksum    string
	Tag            string
	ClonedAt       time.Time
}

// Clone duplicates a snapshot entry identified by ref (checksum or tag),
// optionally applying a new tag to the cloned entry.
func Clone(path, ref, newTag string) (*CloneResult, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}
	if ref == "" {
		return nil, errors.New("ref is required")
	}

	snap, err := Load(path)
	if err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}

	var source *Entry
	for i := range snap.Entries {
		e := &snap.Entries[i]
		if e.Checksum == ref || e.Tag == ref {
			source = e
			break
		}
	}
	if source == nil {
		return nil, fmt.Errorf("ref %q not found in snapshot", ref)
	}

	cloned := Entry{
		Checksum:  source.Checksum + "-clone-" + fmt.Sprintf("%d", time.Now().UnixNano()),
		Keys:      append([]string(nil), source.Keys...),
		CreatedAt: time.Now().UTC(),
		Tag:       newTag,
		Note:      fmt.Sprintf("cloned from %s", source.Checksum),
	}

	snap.Entries = append(snap.Entries, cloned)
	if err := Save(path, snap); err != nil {
		return nil, fmt.Errorf("save snapshot: %w", err)
	}

	return &CloneResult{
		SourceChecksum: source.Checksum,
		NewChecksum:    cloned.Checksum,
		Tag:            cloned.Tag,
		ClonedAt:       cloned.CreatedAt,
	}, nil
}
