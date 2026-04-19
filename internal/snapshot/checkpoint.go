package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Checkpoint struct {
	Checksum  string    `json:"checksum"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

type checkpointStore struct {
	Checkpoints []Checkpoint `json:"checkpoints"`
}

func checkpointPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "checkpoints.json")
}

func loadCheckpointStore(path string) (checkpointStore, error) {
	var store checkpointStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveCheckpointStore(path string, store checkpointStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SetCheckpoint records a named checkpoint for a given checksum.
func SetCheckpoint(snapshotPath, checksum, label, createdBy string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if label == "" {
		return errors.New("label is required")
	}
	if createdBy == "" {
		return errors.New("createdBy is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return err
	}
	entry := findEntry(snap, checksum)
	if entry == nil {
		return errors.New("checksum not found in snapshot")
	}

	cp := checkpointPath(snapshotPath)
	store, err := loadCheckpointStore(cp)
	if err != nil {
		return err
	}
	store.Checkpoints = append(store.Checkpoints, Checkpoint{
		Checksum:  checksum,
		Label:     label,
		CreatedAt: time.Now().UTC(),
		CreatedBy: createdBy,
	})
	return saveCheckpointStore(cp, store)
}

// GetCheckpoints returns all checkpoints for a given checksum.
func GetCheckpoints(snapshotPath, checksum string) ([]Checkpoint, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	cp := checkpointPath(snapshotPath)
	store, err := loadCheckpointStore(cp)
	if err != nil {
		return nil, err
	}
	var result []Checkpoint
	for _, c := range store.Checkpoints {
		if c.Checksum == checksum {
			result = append(result, c)
		}
	}
	return result, nil
}

func findEntry(snap []Entry, checksum string) *Entry {
	for i := range snap {
		if snap[i].Checksum == checksum {
			return &snap[i]
		}
	}
	return nil
}
