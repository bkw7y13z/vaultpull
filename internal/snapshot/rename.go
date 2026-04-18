package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// RenameResult holds the outcome of a tag rename operation.
type RenameResult struct {
	Checksum string
	OldTag   string
	NewTag   string
	RenamedAt time.Time
}

// RenameTag finds an entry by its existing tag and replaces it with newTag.
// Returns an error if the path, oldTag, or newTag are empty, or if no entry
// carries oldTag.
func RenameTag(path, oldTag, newTag string) (*RenameResult, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}
	if oldTag == "" {
		return nil, errors.New("old tag is required")
	}
	if newTag == "" {
		return nil, errors.New("new tag is required")
	}

	snap, err := Load(path)
	if err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}

	var matched *Entry
	for i := range snap.Entries {
		if snap.Entries[i].Tag == oldTag {
			matched = &snap.Entries[i]
			break
		}
	}
	if matched == nil {
		return nil, fmt.Errorf("no entry found with tag %q", oldTag)
	}

	old := matched.Tag
	matched.Tag = newTag

	if err := Save(path, snap); err != nil {
		return nil, fmt.Errorf("save snapshot: %w", err)
	}

	return &RenameResult{
		Checksum:  matched.Checksum,
		OldTag:    old,
		NewTag:    newTag,
		RenamedAt: time.Now().UTC(),
	}, nil
}
