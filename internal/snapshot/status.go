package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type StatusEntry struct {
	Checksum string `json:"checksum"`
	State    string `json:"state"` // "active", "deprecated", "archived"
	SetBy    string `json:"set_by"`
	Reason   string `json:"reason"`
}

type statusStore struct {
	Entries map[string]StatusEntry `json:"entries"`
}

func statusPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "status.json")
}

func loadStatusStore(path string) (statusStore, error) {
	var store statusStore
	store.Entries = make(map[string]StatusEntry)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveStatusStore(path string, store statusStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SetStatus(snapshotPath, checksum, state, setBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	valid := map[string]bool{"active": true, "deprecated": true, "archived": true}
	if !valid[state] {
		return errors.New("state must be one of: active, deprecated, archived")
	}
	if setBy == "" {
		return errors.New("set_by is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return err
	}
	found := false
	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			found = true
			break
		}
	}
	if !found {
		return errors.New("checksum not found in snapshot")
	}

	store, err := loadStatusStore(statusPath(snapshotPath))
	if err != nil {
		return err
	}
	store.Entries[checksum] = StatusEntry{
		Checksum: checksum,
		State:    state,
		SetBy:    setBy,
		Reason:   reason,
	}
	return saveStatusStore(statusPath(snapshotPath), store)
}

func GetStatus(snapshotPath, checksum string) (StatusEntry, error) {
	if snapshotPath == "" {
		return StatusEntry{}, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return StatusEntry{}, errors.New("checksum is required")
	}
	store, err := loadStatusStore(statusPath(snapshotPath))
	if err != nil {
		return StatusEntry{}, err
	}
	entry, ok := store.Entries[checksum]
	if !ok {
		return StatusEntry{}, errors.New("no status found for checksum")
	}
	return entry, nil
}
