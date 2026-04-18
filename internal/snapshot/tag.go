package snapshot

import (
	"errors"
	"fmt"
	"os"
	"time"

	"encoding/json"
)

// Tag assigns a human-readable label to a snapshot entry by its checksum.
func Tag(snapshotPath, checksum, label string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path must not be empty")
	}
	if checksum == "" {
		return errors.New("checksum must not be empty")
	}
	if label == "" {
		return errors.New("label must not be empty")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	found := false
	for i, e := range entries {
		if e.Checksum == checksum {
			entries[i].Tag = label
			entries[i].TaggedAt = time.Now().UTC()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no snapshot entry found with checksum %q", checksum)
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("marshalling snapshot: %w", err)
	}

	if err := os.WriteFile(snapshotPath, data, 0600); err != nil {
		return fmt.Errorf("writing snapshot: %w", err)
	}

	return nil
}

// FindByTag returns the first snapshot entry matching the given label.
func FindByTag(snapshotPath, label string) (*Entry, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path must not be empty")
	}
	if label == "" {
		return nil, errors.New("label must not be empty")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot: %w", err)
	}

	for _, e := range entries {
		if e.Tag == label {
			copy := e
			return &copy, nil
		}
	}

	return nil, fmt.Errorf("no snapshot entry found with tag %q", label)
}

// RemoveTag clears the tag and tagged timestamp from the snapshot entry
// identified by the given checksum.
func RemoveTag(snapshotPath, checksum string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path must not be empty")
	}
	if checksum == "" {
		return errors.New("checksum must not be empty")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}

	found := false
	for i, e := range entries {
		if e.Checksum == checksum {
			entries[i].Tag = ""
			entries[i].TaggedAt = time.Time{}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("no snapshot entry found with checksum %q", checksum)
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("marshalling snapshot: %w", err)
	}

	if err := os.WriteFile(snapshotPath, data, 0600); err != nil {
		return fmt.Errorf("writing snapshot: %w", err)
	}

	return nil
}
