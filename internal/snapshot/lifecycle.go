package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type LifecycleState string

const (
	LifecycleActive     LifecycleState = "active"
	LifecycleDeprecated LifecycleState = "deprecated"
	LifecycleRetired    LifecycleState = "retired"
)

type LifecycleEntry struct {
	Checksum   string         `json:"checksum"`
	State      LifecycleState `json:"state"`
	ChangedBy  string         `json:"changed_by"`
	Reason     string         `json:"reason"`
	ChangedAt  time.Time      `json:"changed_at"`
}

type lifecycleStore struct {
	Entries []LifecycleEntry `json:"entries"`
}

func lifecyclePath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "lifecycle.json")
}

func loadLifecycleStore(path string) (lifecycleStore, error) {
	var store lifecycleStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveLifecycleStore(path string, store lifecycleStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SetLifecycle(snapshotPath, checksum string, state LifecycleState, changedBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if changedBy == "" {
		return errors.New("changed_by is required")
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

	p := lifecyclePath(snapshotPath)
	store, err := loadLifecycleStore(p)
	if err != nil {
		return err
	}
	for i, e := range store.Entries {
		if e.Checksum == checksum {
			store.Entries[i] = LifecycleEntry{Checksum: checksum, State: state, ChangedBy: changedBy, Reason: reason, ChangedAt: time.Now().UTC()}
			return saveLifecycleStore(p, store)
		}
	}
	store.Entries = append(store.Entries, LifecycleEntry{Checksum: checksum, State: state, ChangedBy: changedBy, Reason: reason, ChangedAt: time.Now().UTC()})
	return saveLifecycleStore(p, store)
}

func GetLifecycle(snapshotPath, checksum string) (LifecycleEntry, bool, error) {
	if snapshotPath == "" {
		return LifecycleEntry{}, false, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return LifecycleEntry{}, false, errors.New("checksum is required")
	}
	p := lifecyclePath(snapshotPath)
	store, err := loadLifecycleStore(p)
	if err != nil {
		return LifecycleEntry{}, false, err
	}
	for _, e := range store.Entries {
		if e.Checksum == checksum {
			return e, true, nil
		}
	}
	return LifecycleEntry{}, false, nil
}
