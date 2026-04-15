package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single snapshot record of secrets pulled from Vault.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Namespace string            `json:"namespace"`
	Keys      []string          `json:"keys"`
	Checksum  string            `json:"checksum"`
}

// Snapshot holds a collection of snapshot entries persisted to disk.
type Snapshot struct {
	Entries []Entry `json:"entries"`
}

// Load reads a snapshot file from the given path.
// Returns an empty Snapshot if the file does not exist.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Snapshot{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read file: %w", err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &s, nil
}

// Save persists the snapshot to the given path, creating or overwriting it.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("snapshot: write file: %w", err)
	}
	return nil
}

// Add appends a new entry to the snapshot, keeping at most maxEntries.
func (s *Snapshot) Add(entry Entry, maxEntries int) {
	s.Entries = append(s.Entries, entry)
	if maxEntries > 0 && len(s.Entries) > maxEntries {
		s.Entries = s.Entries[len(s.Entries)-maxEntries:]
	}
}

// Latest returns the most recent snapshot entry, or nil if none exist.
func (s *Snapshot) Latest() *Entry {
	if len(s.Entries) == 0 {
		return nil
	}
	e := s.Entries[len(s.Entries)-1]
	return &e
}
