package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type LabelStore struct {
	Labels map[string]string `json:"labels"` // checksum -> label
}

func labelPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "labels.json")
}

func loadLabelStore(path string) (*LabelStore, error) {
	store := &LabelStore{Labels: map[string]string{}}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	return store, json.Unmarshal(data, store)
}

func saveLabelStore(path string, store *LabelStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Label sets a human-readable label for a snapshot entry by checksum.
func Label(snapshotPath, checksum, label string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if label == "" {
		return errors.New("label is required")
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

	store, err := loadLabelStore(labelPath(snapshotPath))
	if err != nil {
		return err
	}
	store.Labels[checksum] = label
	return saveLabelStore(labelPath(snapshotPath), store)
}

// GetLabel returns the label for a given checksum.
func GetLabel(snapshotPath, checksum string) (string, error) {
	if snapshotPath == "" {
		return "", errors.New("snapshot path is required")
	}
	if checksum == "" {
		return "", errors.New("checksum is required")
	}
	store, err := loadLabelStore(labelPath(snapshotPath))
	if err != nil {
		return "", err
	}
	lbl, ok := store.Labels[checksum]
	if !ok {
		return "", errors.New("no label found for checksum")
	}
	return lbl, nil
}
