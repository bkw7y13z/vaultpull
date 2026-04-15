package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const maxEntries = 100

// Entry represents a single snapshot record.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Checksum  string            `json:"checksum"`
	Keys      []string          `json:"keys"`
	Secrets   map[string]string `json:"secrets"`
	Tag       string            `json:"tag,omitempty"`
	TaggedAt  time.Time         `json:"tagged_at,omitempty"`
}

// Load reads all snapshot entries from the given file path.
func Load(path string) ([]Entry, error) {
	if path == "" {
		return nil, errors.New("snapshot path must not be empty")
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading snapshot file: %w", err)
	}

	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("parsing snapshot file: %w", err)
	}

	return entries, nil
}

// Save persists entries to the snapshot file, capping at maxEntries.
func Save(path string, entries []Entry) error {
	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("marshalling snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing snapshot: %w", err)
	}

	return nil
}

// Add appends a new entry and saves the snapshot.
func Add(path string, entry Entry) error {
	entries, err := Load(path)
	if err != nil {
		return err
	}
	entries = append(entries, entry)
	return Save(path, entries)
}

// Latest returns the most recent snapshot entry, if any.
func Latest(path string) (*Entry, error) {
	entries, err := Load(path)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	e := entries[len(entries)-1]
	return &e, nil
}
