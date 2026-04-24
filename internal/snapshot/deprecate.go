package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type DeprecationRecord struct {
	Checksum    string    `json:"checksum"`
	Reason      string    `json:"reason"`
	DeprecatedBy string   `json:"deprecated_by"`
	DeprecatedAt time.Time `json:"deprecated_at"`
	Suggest     string    `json:"suggest,omitempty"`
}

type deprecateStore struct {
	Records []DeprecationRecord `json:"records"`
}

func deprecatePath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "deprecations.json")
}

func loadDeprecateStore(path string) (deprecateStore, error) {
	var store deprecateStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveDeprecateStore(path string, store deprecateStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Deprecate marks a snapshot entry as deprecated with a reason and optional suggestion.
func Deprecate(snapshotPath, checksum, reason, deprecatedBy, suggest string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if reason == "" {
		return errors.New("reason is required")
	}
	if deprecatedBy == "" {
		return errors.New("deprecated_by is required")
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

	store, err := loadDeprecateStore(deprecatePath(snapshotPath))
	if err != nil {
		return err
	}
	store.Records = append(store.Records, DeprecationRecord{
		Checksum:     checksum,
		Reason:       reason,
		DeprecatedBy: deprecatedBy,
		DeprecatedAt: time.Now().UTC(),
		Suggest:      suggest,
	})
	return saveDeprecateStore(deprecatePath(snapshotPath), store)
}

// GetDeprecation returns the deprecation record for a given checksum, if any.
func GetDeprecation(snapshotPath, checksum string) (*DeprecationRecord, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	store, err := loadDeprecateStore(deprecatePath(snapshotPath))
	if err != nil {
		return nil, err
	}
	for i := len(store.Records) - 1; i >= 0; i-- {
		if store.Records[i].Checksum == checksum {
			rec := store.Records[i]
			return &rec, nil
		}
	}
	return nil, nil
}
