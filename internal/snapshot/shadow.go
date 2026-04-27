package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ShadowEntry records a "shadow" copy of a snapshot entry's metadata
// at a point in time, allowing before/after comparison without full cloning.
type ShadowEntry struct {
	Checksum  string            `json:"checksum"`
	Keys      []string          `json:"keys"`
	Tag       string            `json:"tag,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	Meta      map[string]string `json:"meta,omitempty"`
}

type shadowStore struct {
	Shadows []ShadowEntry `json:"shadows"`
}

func shadowPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "shadow_store.json")
}

func loadShadowStore(path string) (shadowStore, error) {
	var store shadowStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveShadowStore(path string, store shadowStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Shadow captures a shadow copy of the snapshot entry identified by checksum.
func Shadow(snapshotPath, checksum string, meta map[string]string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return err
	}

	var found *Entry
	for i := range snap.Entries {
		if snap.Entries[i].Checksum == checksum {
			found = &snap.Entries[i]
			break
		}
	}
	if found == nil {
		return errors.New("checksum not found in snapshot")
	}

	sp := shadowPath(snapshotPath)
	store, err := loadShadowStore(sp)
	if err != nil {
		return err
	}

	entry := ShadowEntry{
		Checksum:  found.Checksum,
		Keys:      KeysFromSecrets(found.Secrets),
		Tag:       found.Tag,
		CreatedAt: time.Now().UTC(),
		Meta:      meta,
	}
	store.Shadows = append(store.Shadows, entry)
	return saveShadowStore(sp, store)
}

// GetShadows returns all shadow entries recorded for the given checksum.
func GetShadows(snapshotPath, checksum string) ([]ShadowEntry, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	store, err := loadShadowStore(shadowPath(snapshotPath))
	if err != nil {
		return nil, err
	}
	var results []ShadowEntry
	for _, s := range store.Shadows {
		if s.Checksum == checksum {
			results = append(results, s)
		}
	}
	return results, nil
}
