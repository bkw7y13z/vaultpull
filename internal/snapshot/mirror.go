package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MirrorEntry records a mirrored snapshot destination.
type MirrorEntry struct {
	Checksum  string    `json:"checksum"`
	DestPath  string    `json:"dest_path"`
	MirroredAt time.Time `json:"mirrored_at"`
	MirroredBy string    `json:"mirrored_by"`
}

func mirrorPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "mirror_index.json")
}

func loadMirrorStore(path string) ([]MirrorEntry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []MirrorEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []MirrorEntry
	return entries, json.Unmarshal(data, &entries)
}

func saveMirrorStore(path string, entries []MirrorEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Mirror copies the snapshot file identified by checksum to destPath and records the operation.
func Mirror(snapshotPath, checksum, destPath, mirroredBy string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if destPath == "" {
		return errors.New("dest path is required")
	}
	if mirroredBy == "" {
		return errors.New("mirrored_by is required")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	var found *Entry
	for i := range entries {
		if entries[i].Checksum == checksum {
			found = &entries[i]
			break
		}
	}
	if found == nil {
		return fmt.Errorf("checksum %q not found", checksum)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("create dest dir: %w", err)
	}

	data, err := os.ReadFile(snapshotPath)
	if err != nil {
		return fmt.Errorf("read snapshot: %w", err)
	}
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("write mirror: %w", err)
	}

	store, err := loadMirrorStore(mirrorPath(snapshotPath))
	if err != nil {
		return err
	}
	store = append(store, MirrorEntry{
		Checksum:   checksum,
		DestPath:   destPath,
		MirroredAt: time.Now().UTC(),
		MirroredBy: mirroredBy,
	})
	return saveMirrorStore(mirrorPath(snapshotPath), store)
}
