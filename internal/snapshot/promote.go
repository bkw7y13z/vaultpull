package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type PromotionRecord struct {
	Checksum    string    `json:"checksum"`
	FromEnv     string    `json:"from_env"`
	ToEnv       string    `json:"to_env"`
	PromotedBy  string    `json:"promoted_by"`
	PromotedAt  time.Time `json:"promoted_at"`
	Note        string    `json:"note,omitempty"`
}

type PromoteStore struct {
	Records []PromotionRecord `json:"records"`
}

func promotePath(snapshotPath string) string {
	ext := filepath.Ext(snapshotPath)
	return snapshotPath[:len(snapshotPath)-len(ext)] + ".promotions.json"
}

func loadPromoteStore(path string) (PromoteStore, error) {
	var store PromoteStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func savePromoteStore(path string, store PromoteStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Promote(snapshotPath, checksum, fromEnv, toEnv, promotedBy, note string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if fromEnv == "" {
		return errors.New("from_env is required")
	}
	if toEnv == "" {
		return errors.New("to_env is required")
	}
	if promotedBy == "" {
		return errors.New("promoted_by is required")
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

	pp := promotePath(snapshotPath)
	store, err := loadPromoteStore(pp)
	if err != nil {
		return err
	}
	store.Records = append(store.Records, PromotionRecord{
		Checksum:   checksum,
		FromEnv:    fromEnv,
		ToEnv:      toEnv,
		PromotedBy: promotedBy,
		PromotedAt: time.Now().UTC(),
		Note:       note,
	})
	return savePromoteStore(pp, store)
}

func GetPromotions(snapshotPath, checksum string) ([]PromotionRecord, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	pp := promotePath(snapshotPath)
	store, err := loadPromoteStore(pp)
	if err != nil {
		return nil, err
	}
	var out []PromotionRecord
	for _, r := range store.Records {
		if r.Checksum == checksum {
			out = append(out, r)
		}
	}
	return out, nil
}
