package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Badge struct {
	Checksum string `json:"checksum"`
	Label    string `json:"label"`
	Status   string `json:"status"` // ok, warning, error
	Message  string `json:"message"`
}

type BadgeStore struct {
	Badges map[string]Badge `json:"badges"`
}

func badgePath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "badges.json")
}

func loadBadgeStore(path string) (BadgeStore, error) {
	var store BadgeStore
	store.Badges = make(map[string]Badge)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveBadgeStore(path string, store BadgeStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func SetBadge(snapshotPath, checksum, label, status, message string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if label == "" {
		return errors.New("label is required")
	}
	validStatuses := map[string]bool{"ok": true, "warning": true, "error": true}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status %q: must be ok, warning, or error", status)
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

	store, err := loadBadgeStore(badgePath(snapshotPath))
	if err != nil {
		return err
	}
	store.Badges[checksum] = Badge{Checksum: checksum, Label: label, Status: status, Message: message}
	return saveBadgeStore(badgePath(snapshotPath), store)
}

func GetBadge(snapshotPath, checksum string) (Badge, bool, error) {
	if snapshotPath == "" {
		return Badge{}, false, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return Badge{}, false, errors.New("checksum is required")
	}
	store, err := loadBadgeStore(badgePath(snapshotPath))
	if err != nil {
		return Badge{}, false, err
	}
	b, ok := store.Badges[checksum]
	return b, ok, nil
}
