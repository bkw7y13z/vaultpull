package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const maxEntries = 100

// Entry represents a single captured snapshot.
type Entry struct {
	Checksum   string            `json:"checksum"`
	Keys       []string          `json:"keys"`
	CapturedAt time.Time         `json:"captured_at"`
	Tag        string            `json:"tag,omitempty"`
	Note       string            `json:"note,omitempty"`
	Meta       map[string]string `json:"meta,omitempty"`
}

// Snapshot holds a list of snapshot entries.
type Snapshot struct {
	Entries []Entry `json:"entries"`
}

// Load reads a snapshot from disk.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return &Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &snap, nil
}

// Save writes a snapshot to disk.
func Save(path string, snap *Snapshot) error {
	if snap == nil {
		return errors.New("snapshot is nil")
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	return os.WriteFile(path, data, 0600)
}

// Add appends an entry, capping at maxEntries.
func Add(path string, entry Entry) error {
	snap, err := Load(path)
	ifNotExist(err) {
		return fmt.Errorf("load: %w", err)
	}
	snap.Entries = append(snap.Entries, entry)
	if len(snap.Entries) > maxEntries {
		snap.Entries = snap.Entries[len(snap.Entries)-maxEntries:]
	}
	return Save(path, snap)
}

// Latest returns the most recent entry, or nil if empty.
func Latest(path string) (*Entry, error) {
	snap, err := Load(path)
	if err != nil {
		return nil, err
	}
	if len(snap.Entries) == 0 {
		return nil, nil
	}
	e := snap.Entries[len(snap.Entries)-1]
	return &e, nil
}

// FindByTag returns the most recent entry matching the given tag, or nil if none found.
func FindByTag(path string, tag string) (*Entry, error) {
	snap, err := Load(path)
	if err != nil {
		return nil, err
	}
	for i := len(snap.Entries) - 1; i >= 0; i-- {
		if snap.Entries[i].Tag == tag {
			e := snap.Entries[i]
			return &e, nil
		}
	}
	return nil, nil
}
