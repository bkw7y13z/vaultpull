package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type LockEntry struct {
	Checksum  string    `json:"checksum"`
	LockedBy  string    `json:"locked_by"`
	Reason    string    `json:"reason"`
	LockedAt  time.Time `json:"locked_at"`
}

type lockStore struct {
	Locks []LockEntry `json:"locks"`
}

func lockPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "locks.json")
}

func loadLockStore(path string) (lockStore, error) {
	var store lockStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveLockStore(path string, store lockStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Lock(snapshotPath, checksum, lockedBy, reason string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if lockedBy == "" {
		return errors.New("locked_by is required")
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

	lp := lockPath(snapshotPath)
	store, err := loadLockStore(lp)
	if err != nil {
		return err
	}
	for _, l := range store.Locks {
		if l.Checksum == checksum {
			return errors.New("entry is already locked")
		}
	}
	store.Locks = append(store.Locks, LockEntry{
		Checksum: checksum,
		LockedBy: lockedBy,
		Reason:   reason,
		LockedAt: time.Now().UTC(),
	})
	return saveLockStore(lp, store)
}

func IsLocked(snapshotPath, checksum string) (bool, *LockEntry, error) {
	if snapshotPath == "" {
		return false, nil, errors.New("snapshot path is required")
	}
	lp := lockPath(snapshotPath)
	store, err := loadLockStore(lp)
	if err != nil {
		return false, nil, err
	}
	for i, l := range store.Locks {
		if l.Checksum == checksum {
			return true, &store.Locks[i], nil
		}
	}
	return false, nil, nil
}

func Unlock(snapshotPath, checksum string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	lp := lockPath(snapshotPath)
	store, err := loadLockStore(lp)
	if err != nil {
		return err
	}
	origLen := len(store.Locks)
	filtered := store.Locks[:0]
	for _, l := range store.Locks {
		if l.Checksum != checksum {
			filtered = append(filtered, l)
		}
	}
	if len(filtered) == origLen {
		return errors.New("lock not found")
	}
	store.Locks = filtered
	return saveLockStore(lp, store)
}
