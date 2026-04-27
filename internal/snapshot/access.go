package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type AccessEntry struct {
	Checksum  string    `json:"checksum"`
	AccessedBy string   `json:"accessed_by"`
	Action    string    `json:"action"`
	At        time.Time `json:"at"`
	Reason    string    `json:"reason,omitempty"`
}

type accessStore struct {
	Entries []AccessEntry `json:"entries"`
}

func accessPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "access_log.json")
}

func loadAccessStore(path string) (accessStore, error) {
	var store accessStore
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveAccessStore(path string, store accessStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RecordAccess logs an access event for a snapshot entry identified by checksum.
func RecordAccess(snapshotPath, checksum, accessedBy, action, reason string) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return fmt.Errorf("checksum is required")
	}
	if accessedBy == "" {
		return fmt.Errorf("accessed_by is required")
	}
	if action == "" {
		return fmt.Errorf("action is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
	}
	if _, ok := snap[checksum]; !ok {
		return fmt.Errorf("checksum %q not found in snapshot", checksum)
	}

	ap := accessPath(snapshotPath)
	store, err := loadAccessStore(ap)
	if err != nil {
		return fmt.Errorf("loading access store: %w", err)
	}

	store.Entries = append(store.Entries, AccessEntry{
		Checksum:   checksum,
		AccessedBy: accessedBy,
		Action:     action,
		At:         time.Now().UTC(),
		Reason:     reason,
	})

	return saveAccessStore(ap, store)
}

// GetAccessLog returns all access entries for a given checksum.
func GetAccessLog(snapshotPath, checksum string) ([]AccessEntry, error) {
	if snapshotPath == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return nil, fmt.Errorf("checksum is required")
	}

	ap := accessPath(snapshotPath)
	store, err := loadAccessStore(ap)
	if err != nil {
		return nil, err
	}

	var result []AccessEntry
	for _, e := range store.Entries {
		if e.Checksum == checksum {
			result = append(result, e)
		}
	}
	return result, nil
}
