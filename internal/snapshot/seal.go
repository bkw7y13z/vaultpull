package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// SealRecord represents a sealed (locked) snapshot entry that cannot be pruned or modified.
type SealRecord struct {
	Checksum  string    `json:"checksum"`
	SealedAt  time.Time `json:"sealed_at"`
	SealedBy  string    `json:"sealed_by"`
	Reason    string    `json:"reason"`
}

type sealStore struct {
	Seals []SealRecord `json:"seals"`
}

// Seal marks a snapshot entry as sealed, preventing modification or pruning.
func Seal(path, checksum, sealedBy, reason string) error {
	if path == "" {
		return errors.New("seal: path is required")
	}
	if checksum == "" {
		return errors.New("seal: checksum is required")
	}
	if sealedBy == "" {
		return errors.New("seal: sealedBy is required")
	}
	if reason == "" {
		return errors.New("seal: reason is required")
	}

	store, err := loadSealStore(path)
	if err != nil {
		return err
	}

	for _, s := range store.Seals {
		if s.Checksum == checksum {
			return errors.New("seal: entry already sealed")
		}
	}

	store.Seals = append(store.Seals, SealRecord{
		Checksum: checksum,
		SealedAt: time.Now().UTC(),
		SealedBy: sealedBy,
		Reason:   reason,
	})

	return saveSealStore(path, store)
}

// IsSealed returns true if the given checksum is sealed.
func IsSealed(path, checksum string) (bool, error) {
	store, err := loadSealStore(path)
	if err != nil {
		return false, err
	}
	for _, s := range store.Seals {
		if s.Checksum == checksum {
			return true, nil
		}
	}
	return false, nil
}

// GetSeal returns the SealRecord for a checksum, or an error if not found.
func GetSeal(path, checksum string) (SealRecord, error) {
	store, err := loadSealStore(path)
	if err != nil {
		return SealRecord{}, err
	}
	for _, s := range store.Seals {
		if s.Checksum == checksum {
			return s, nil
		}
	}
	return SealRecord{}, errors.New("seal: record not found")
}

func loadSealStore(path string) (sealStore, error) {
	var store sealStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveSealStore(path string, store sealStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
