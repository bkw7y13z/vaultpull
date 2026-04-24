package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DependencyEdge represents a directional dependency between two snapshot entries.
type DependencyEdge struct {
	FromChecksum string    `json:"from_checksum"`
	ToChecksum   string    `json:"to_checksum"`
	Reason       string    `json:"reason"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
}

// DependencyStore holds all recorded dependency edges for a snapshot file.
type DependencyStore struct {
	Edges []DependencyEdge `json:"edges"`
}

func dependencyPath(snapshotPath string) string {
	ext := filepath.Ext(snapshotPath)
	base := snapshotPath[:len(snapshotPath)-len(ext)]
	return base + ".dependencies.json"
}

func loadDependencyStore(path string) (DependencyStore, error) {
	var store DependencyStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, fmt.Errorf("read dependency store: %w", err)
	}
	if err := json.Unmarshal(data, &store); err != nil {
		return store, fmt.Errorf("parse dependency store: %w", err)
	}
	return store, nil
}

func saveDependencyStore(path string, store DependencyStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal dependency store: %w", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write dependency store: %w", err)
	}
	return nil
}

// AddDependency records a dependency edge from one snapshot entry to another.
// Both checksums must exist in the snapshot file.
func AddDependency(snapshotPath, fromChecksum, toChecksum, reason, createdBy string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if fromChecksum == "" {
		return errors.New("from_checksum is required")
	}
	if toChecksum == "" {
		return errors.New("to_checksum is required")
	}
	if reason == "" {
		return errors.New("reason is required")
	}
	if createdBy == "" {
		return errors.New("created_by is required")
	}
	if fromChecksum == toChecksum {
		return errors.New("from_checksum and to_checksum must differ")
	}

	// Validate both checksums exist in the snapshot.
	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	hasFrom, hasTo := false, false
	for _, e := range snap {
		if e.Checksum == fromChecksum {
			hasFrom = true
		}
		if e.Checksum == toChecksum {
			hasTo = true
		}
	}
	if !hasFrom {
		return fmt.Errorf("from_checksum %q not found in snapshot", fromChecksum)
	}
	if !hasTo {
		return fmt.Errorf("to_checksum %q not found in snapshot", toChecksum)
	}

	depPath := dependencyPath(snapshotPath)
	store, err := loadDependencyStore(depPath)
	if err != nil {
		return err
	}

	edge := DependencyEdge{
		FromChecksum: fromChecksum,
		ToChecksum:   toChecksum,
		Reason:       reason,
		CreatedBy:    createdBy,
		CreatedAt:    time.Now().UTC(),
	}
	store.Edges = append(store.Edges, edge)
	return saveDependencyStore(depPath, store)
}

// GetDependencies returns all dependency edges originating from the given checksum.
func GetDependencies(snapshotPath, fromChecksum string) ([]DependencyEdge, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if fromChecksum == "" {
		return nil, errors.New("from_checksum is required")
	}

	depPath := dependencyPath(snapshotPath)
	store, err := loadDependencyStore(depPath)
	if err != nil {
		return nil, err
	}

	var result []DependencyEdge
	for _, e := range store.Edges {
		if e.FromChecksum == fromChecksum {
			result = append(result, e)
		}
	}
	return result, nil
}
