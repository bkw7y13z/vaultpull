package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type EvictionRecord struct {
	Checksum  string    `json:"checksum"`
	EvictedBy string    `json:"evicted_by"`
	Reason    string    `json:"reason"`
	EvictedAt time.Time `json:"evicted_at"`
}

type evictStore struct {
	Evictions []EvictionRecord `json:"evictions"`
}

func evictPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "evictions.json")
}

func loadEvictStore(path string) (evictStore, error) {
	var store evictStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveEvictStore(path string, store evictStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// Evict marks a snapshot entry as evicted and removes it from the snapshot file.
func Evict(snapshotPath, checksum, evictedBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if evictedBy == "" {
		return errors.New("evicted_by is required")
	}
	if reason == "" {
		return errors.New("reason is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	found := false
	filtered := snap.Entries[:0]
	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			found = true
			continue
		}
		filtered = append(filtered, e)
	}
	if !found {
		return fmt.Errorf("checksum %q not found in snapshot", checksum)
	}
	snap.Entries = filtered
	if err := Save(snapshotPath, snap); err != nil {
		return fmt.Errorf("save snapshot: %w", err)
	}

	ep := evictPath(snapshotPath)
	store, err := loadEvictStore(ep)
	if err != nil {
		return fmt.Errorf("load evict store: %w", err)
	}
	store.Evictions = append(store.Evictions, EvictionRecord{
		Checksum:  checksum,
		EvictedBy: evictedBy,
		Reason:    reason,
		EvictedAt: time.Now().UTC(),
	})
	return saveEvictStore(ep, store)
}

// GetEvictions returns all eviction records for the given snapshot.
func GetEvictions(snapshotPath string) ([]EvictionRecord, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	ep := evictPath(snapshotPath)
	store, err := loadEvictStore(ep)
	if err != nil {
		return nil, err
	}
	return store.Evictions, nil
}
