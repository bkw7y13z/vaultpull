package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type QuotaPolicy struct {
	Checksum  string    `json:"checksum"`
	MaxKeys   int       `json:"max_keys"`
	MaxSizeKB int       `json:"max_size_kb"`
	SetBy     string    `json:"set_by"`
	CreatedAt time.Time `json:"created_at"`
}

type QuotaStore struct {
	Policies map[string]QuotaPolicy `json:"policies"`
}

func quotaPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "quota.json")
}

func loadQuotaStore(path string) (QuotaStore, error) {
	store := QuotaStore{Policies: make(map[string]QuotaPolicy)}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveQuotaStore(path string, store QuotaStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SetQuota(snapshotPath, checksum, setBy string, maxKeys, maxSizeKB int) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if setBy == "" {
		return errors.New("set_by is required")
	}
	if maxKeys <= 0 && maxSizeKB <= 0 {
		return errors.New("at least one of max_keys or max_size_kb must be positive")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("loading snapshot: %w", err)
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

	qPath := quotaPath(snapshotPath)
	store, err := loadQuotaStore(qPath)
	if err != nil {
		return err
	}
	store.Policies[checksum] = QuotaPolicy{
		Checksum:  checksum,
		MaxKeys:   maxKeys,
		MaxSizeKB: maxSizeKB,
		SetBy:     setBy,
		CreatedAt: time.Now().UTC(),
	}
	return saveQuotaStore(qPath, store)
}

func GetQuota(snapshotPath, checksum string) (QuotaPolicy, bool, error) {
	if snapshotPath == "" {
		return QuotaPolicy{}, false, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return QuotaPolicy{}, false, errors.New("checksum is required")
	}
	qPath := quotaPath(snapshotPath)
	store, err := loadQuotaStore(qPath)
	if err != nil {
		return QuotaPolicy{}, false, err
	}
	p, ok := store.Policies[checksum]
	return p, ok, nil
}
