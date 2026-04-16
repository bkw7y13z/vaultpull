package snapshot

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// Annotation holds a user-defined note attached to a snapshot entry.
type Annotation struct {
	Checksum  string    `json:"checksum"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// Annotate attaches a note to the snapshot entry matching checksum.
func Annotate(snapshotPath, checksum, note string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if note == "" {
		return errors.New("note is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("load snapshot: %w", err)
	}

	found := false
	for i := range snap.Entries {
		if snap.Entries[i].Checksum == checksum {
			snap.Entries[i].Note = note
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("checksum %q not found in snapshot", checksum)
	}

	return Save(snapshotPath, snap)
}

// GetAnnotation returns the note for a given checksum, or empty string if none.
func GetAnnotation(snapshotPath, checksum string) (string, error) {
	if snapshotPath == "" {
		return "", errors.New("snapshot path is required")
	}
	if checksum == "" {
		return "", errors.New("checksum is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return "", fmt.Errorf("load snapshot: %w", err)
	}

	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			return e.Note, nil
		}
	}
	return "", fmt.Errorf("checksum %q not found", checksum)
}
