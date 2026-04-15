package snapshot

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// RestoreResult holds the outcome of a restore operation.
type RestoreResult struct {
	Checksum  string
	Tag       string
	Timestamp time.Time
	Keys      []string
	Secrets   map[string]string
}

// Restore finds a snapshot entry by checksum or tag and returns its secrets
// as a map suitable for writing to an env file.
func Restore(snapshotPath, ref string) (*RestoreResult, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path must not be empty")
	}
	if ref == "" {
		return nil, errors.New("ref (checksum or tag) must not be empty")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot: %w", err)
	}
	if len(entries) == 0 {
		return nil, errors.New("snapshot file is empty")
	}

	// Try to find by tag first, then by checksum prefix.
	var matched *Entry
	for i := range entries {
		e := &entries[i]
		if e.Tag == ref || e.Checksum == ref || (len(ref) >= 7 && len(e.Checksum) >= len(ref) && e.Checksum[:len(ref)] == ref) {
			matched = e
			break
		}
	}
	if matched == nil {
		return nil, fmt.Errorf("no snapshot entry found for ref %q", ref)
	}

	keys := KeysFromSecrets(matched.Secrets)
	sort.Strings(keys)

	return &RestoreResult{
		Checksum:  matched.Checksum,
		Tag:       matched.Tag,
		Timestamp: matched.Timestamp,
		Keys:      keys,
		Secrets:   matched.Secrets,
	}, nil
}
