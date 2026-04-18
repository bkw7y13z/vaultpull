package snapshot

import (
	"errors"
	"time"
)

// ExpireEntry holds expiration metadata for a snapshot entry.
type ExpireEntry struct {
	Checksum  string    `json:"checksum"`
	ExpiresAt time.Time `json:"expires_at"`
	SetBy     string    `json:"set_by"`
}

// ExpireStore holds all expiration records.
type ExpireStore struct {
	Entries []ExpireEntry `json:"entries"`
}

func expirePath(snapshotPath string) string {
	return snapshotPath + ".expire.json"
}

func loadExpireStore(path string) (ExpireStore, error) {
	var store ExpireStore
	err := loadJSON(expirePath(path), &store)
	if err != nil {
		return ExpireStore{}, err
	}
	return store, nil
}

func saveExpireStore(path string, store ExpireStore) error {
	return saveJSON(expirePath(path), store)
}

// SetExpiry sets an expiration time on the entry matching checksum.
func SetExpiry(snapshotPath, checksum, setBy string, expiresAt time.Time) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if setBy == "" {
		return errors.New("set_by is required")
	}
	if expiresAt.IsZero() {
		return errors.New("expires_at is required")
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

	store, _ := loadExpireStore(snapshotPath)
	for i, e := range store.Entries {
		if e.Checksum == checksum {
			store.Entries[i].ExpiresAt = expiresAt
			store.Entries[i].SetBy = setBy
			return saveExpireStore(snapshotPath, store)
		}
	}
	store.Entries = append(store.Entries, ExpireEntry{
		Checksum:  checksum,
		ExpiresAt: expiresAt,
		SetBy:     setBy,
	})
	return saveExpireStore(snapshotPath, store)
}

// GetExpiry returns the expiration entry for a checksum, or nil if not set.
func GetExpiry(snapshotPath, checksum string) (*ExpireEntry, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	store, err := loadExpireStore(snapshotPath)
	if err != nil {
		return nil, err
	}
	for _, e := range store.Entries {
		if e.Checksum == checksum {
			copy := e
			return &copy, nil
		}
	}
	return nil, nil
}
