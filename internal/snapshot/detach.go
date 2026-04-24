package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type DetachRecord struct {
	Checksum   string    `json:"checksum"`
	DetachedBy string    `json:"detached_by"`
	Reason     string    `json:"reason"`
	DetachedAt time.Time `json:"detached_at"`
}

type detachStore struct {
	Records []DetachRecord `json:"records"`
}

func detachPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "detach.json")
}

func loadDetachStore(path string) (detachStore, error) {
	var store detachStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveDetachStore(path string, store detachStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Detach marks a snapshot entry as detached from its origin source,
// recording who detached it and why.
func Detach(snapshotPath, checksum, detachedBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if detachedBy == "" {
		return errors.New("detached_by is required")
	}
	if reason == "" {
		return errors.New("reason is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	found := false
	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("checksum %q not found in snapshot", checksum)
	}

	store, err := loadDetachStore(detachPath(snapshotPath))
	if err != nil {
		return fmt.Errorf("load detach store: %w", err)
	}

	store.Records = append(store.Records, DetachRecord{
		Checksum:   checksum,
		DetachedBy: detachedBy,
		Reason:     reason,
		DetachedAt: time.Now().UTC(),
	})

	return saveDetachStore(detachPath(snapshotPath), store)
}

// GetDetachment returns the detach record for the given checksum, if any.
func GetDetachment(snapshotPath, checksum string) (*DetachRecord, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}

	store, err := loadDetachStore(detachPath(snapshotPath))
	if err != nil {
		return nil, err
	}

	for i := len(store.Records) - 1; i >= 0; i-- {
		if store.Records[i].Checksum == checksum {
			r := store.Records[i]
			return &r, nil
		}
	}
	return nil, nil
}
