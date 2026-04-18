package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type FreezeStore struct {
	Frozen map[string]FreezeRecord `json:"frozen"`
}

type FreezeRecord struct {
	Checksum string    `json:"checksum"`
	FrozenBy string    `json:"frozen_by"`
	Reason   string    `json:"reason"`
	FrozenAt time.Time `json:"frozen_at"`
}

func freezePath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "freeze_store.json")
}

func loadFreezeStore(path string) (*FreezeStore, error) {
	store := &FreezeStore{Frozen: make(map[string]FreezeRecord)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	return store, json.Unmarshal(data, store)
}

func saveFreezeStore(path string, store *FreezeStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Freeze marks a snapshot entry as frozen, preventing modification.
func Freeze(snapshotPath, checksum, frozenBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if frozenBy == "" {
		return errors.New("frozen_by is required")
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

	fp := freezePath(snapshotPath)
	store, err := loadFreezeStore(fp)
	if err != nil {
		return err
	}
	store.Frozen[checksum] = FreezeRecord{
		Checksum: checksum,
		FrozenBy: frozenBy,
		Reason:   reason,
		FrozenAt: time.Now().UTC(),
	}
	return saveFreezeStore(fp, store)
}

// IsFrozen returns true if the given checksum is frozen.
func IsFrozen(snapshotPath, checksum string) (bool, error) {
	store, err := loadFreezeStore(freezePath(snapshotPath))
	if err != nil {
		return false, err
	}
	_, ok := store.Frozen[checksum]
	return ok, nil
}

// GetFreeze returns the FreezeRecord for a checksum, if present.
func GetFreeze(snapshotPath, checksum string) (FreezeRecord, bool, error) {
	store, err := loadFreezeStore(freezePath(snapshotPath))
	if err != nil {
		return FreezeRecord{}, false, err
	}
	r, ok := store.Frozen[checksum]
	return r, ok, nil
}
