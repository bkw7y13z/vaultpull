package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LineageEdge struct {
	Parent    string    `json:"parent"`
	Child     string    `json:"child"`
	CreatedAt time.Time `json:"created_at"`
	Reason    string    `json:"reason"`
}

type LineageStore struct {
	Edges []LineageEdge `json:"edges"`
}

func lineagePath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "lineage.json")
}

func loadLineageStore(path string) (LineageStore, error) {
	var store LineageStore
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveLineageStore(path string, store LineageStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// AddLineage records a parent-child relationship between two snapshot checksums.
func AddLineage(snapshotPath, parentChecksum, childChecksum, reason string) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if parentChecksum == "" {
		return fmt.Errorf("parent checksum is required")
	}
	if childChecksum == "" {
		return fmt.Errorf("child checksum is required")
	}
	if reason == "" {
		return fmt.Errorf("reason is required")
	}

	store, err := loadLineageStore(lineagePath(snapshotPath))
	if err != nil {
		return fmt.Errorf("load lineage store: %w", err)
	}

	store.Edges = append(store.Edges, LineageEdge{
		Parent:    parentChecksum,
		Child:     childChecksum,
		CreatedAt: time.Now().UTC(),
		Reason:    reason,
	})

	return saveLineageStore(lineagePath(snapshotPath), store)
}

// GetLineage returns all edges where the given checksum is the child (its ancestors).
func GetLineage(snapshotPath, checksum string) ([]LineageEdge, error) {
	if snapshotPath == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return nil, fmt.Errorf("checksum is required")
	}

	store, err := loadLineageStore(lineagePath(snapshotPath))
	if err != nil {
		return nil, fmt.Errorf("load lineage store: %w", err)
	}

	var result []LineageEdge
	for _, e := range store.Edges {
		if e.Child == checksum {
			result = append(result, e)
		}
	}
	return result, nil
}
