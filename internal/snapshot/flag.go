package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type FlagEntry struct {
	Checksum  string    `json:"checksum"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	FlaggedBy string    `json:"flagged_by"`
	Reason    string    `json:"reason"`
	At        time.Time `json:"at"`
}

type FlagStore struct {
	Flags []FlagEntry `json:"flags"`
}

func flagPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "flags.json")
}

func loadFlagStore(path string) (FlagStore, error) {
	var store FlagStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveFlagStore(path string, store FlagStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Flag(snapshotPath, checksum, key, value, flaggedBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if key == "" {
		return errors.New("key is required")
	}
	if flaggedBy == "" {
		return errors.New("flagged_by is required")
	}
	if reason == "" {
		return errors.New("reason is required")
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

	p := flagPath(snapshotPath)
	store, err := loadFlagStore(p)
	if err != nil {
		return err
	}
	store.Flags = append(store.Flags, FlagEntry{
		Checksum:  checksum,
		Key:       key,
		Value:     value,
		FlaggedBy: flaggedBy,
		Reason:    reason,
		At:        time.Now().UTC(),
	})
	return saveFlagStore(p, store)
}

func GetFlags(snapshotPath, checksum string) ([]FlagEntry, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	p := flagPath(snapshotPath)
	store, err := loadFlagStore(p)
	if err != nil {
		return nil, err
	}
	var result []FlagEntry
	for _, f := range store.Flags {
		if checksum == "" || f.Checksum == checksum {
			result = append(result, f)
		}
	}
	return result, nil
}
