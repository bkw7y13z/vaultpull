package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Link struct {
	FromChecksum string    `json:"from_checksum"`
	ToChecksum   string    `json:"to_checksum"`
	Reason       string    `json:"reason"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

type LinkStore struct {
	Links []Link `json:"links"`
}

func linkPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "links.json")
}

func loadLinkStore(path string) (LinkStore, error) {
	var store LinkStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveLinkStore(path string, store LinkStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AddLink creates a directional link between two snapshot checksums.
func AddLink(snapshotPath, fromChecksum, toChecksum, reason, createdBy string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if fromChecksum == "" {
		return errors.New("from checksum is required")
	}
	if toChecksum == "" {
		return errors.New("to checksum is required")
	}
	if reason == "" {
		return errors.New("reason is required")
	}
	if createdBy == "" {
		return errors.New("created_by is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return err
	}
	hasFrom, hasTo := false, false
	for _, e := range snap.Entries {
		if e.Checksum == fromChecksum {
			hasFrom = true
		}
		if e.Checksum == toChecksum {
			hasTo = true
		}
	}
	if !hasFrom {
		return errors.New("from checksum not found in snapshot")
	}
	if !hasTo {
		return errors.New("to checksum not found in snapshot")
	}

	p := linkPath(snapshotPath)
	store, err := loadLinkStore(p)
	if err != nil {
		return err
	}
	store.Links = append(store.Links, Link{
		FromChecksum: fromChecksum,
		ToChecksum:   toChecksum,
		Reason:       reason,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now().UTC(),
	})
	return saveLinkStore(p, store)
}

// GetLinks returns all links originating from a given checksum.
func GetLinks(snapshotPath, fromChecksum string) ([]Link, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if fromChecksum == "" {
		return nil, errors.New("from checksum is required")
	}
	p := linkPath(snapshotPath)
	store, err := loadLinkStore(p)
	if err != nil {
		return nil, err
	}
	var result []Link
	for _, l := range store.Links {
		if l.FromChecksum == fromChecksum {
			result = append(result, l)
		}
	}
	return result, nil
}
