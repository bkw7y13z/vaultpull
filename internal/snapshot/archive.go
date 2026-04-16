package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ArchiveIndex maps checksum to archive file path.
type ArchiveIndex struct {
	Entries []ArchiveRef `json:"entries"`
}

type ArchiveRef struct {
	Checksum  string    `json:"checksum"`
	Tag       string    `json:"tag,omitempty"`
	ArchivedAt time.Time `json:"archived_at"`
	FilePath  string    `json:"file_path"`
}

// Archive moves old snapshot entries to an archive directory and writes an index.
func Archive(snapshotPath, archiveDir string, keepLast int) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if archiveDir == "" {
		return errors.New("archive dir is required")
	}
	if keepLast < 1 {
		return errors.New("keepLast must be at least 1")
	}

	store, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	if len(store.Entries) <= keepLast {
		return nil
	}

	sort.Slice(store.Entries, func(i, j int) bool {
		return store.Entries[i].CreatedAt.Before(store.Entries[j].CreatedAt)
	})

	toArchive := store.Entries[:len(store.Entries)-keepLast]
	keep := store.Entries[len(store.Entries)-keepLast:]

	if err := os.MkdirAll(archiveDir, 0o755); err != nil {
		return fmt.Errorf("create archive dir: %w", err)
	}

	indexPath := filepath.Join(archiveDir, "archive_index.json")
	index := &ArchiveIndex{}
	if data, err := os.ReadFile(indexPath); err == nil {
		_ = json.Unmarshal(data, index)
	}

	for _, e := range toArchive {
		fileName := fmt.Sprintf("%s.json", e.Checksum)
		dest := filepath.Join(archiveDir, fileName)
		data, _ := json.MarshalIndent(e, "", "  ")
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return fmt.Errorf("write archive entry: %w", err)
		}
		index.Entries = append(index.Entries, ArchiveRef{
			Checksum:   e.Checksum,
			Tag:        e.Tag,
			ArchivedAt: time.Now().UTC(),
			FilePath:   dest,
		})
	}

	store.Entries = keep
	if err := Save(snapshotPath, store); err != nil {
		return fmt.Errorf("save pruned snapshot: %w", err)
	}

	idxData, _ := json.MarshalIndent(index, "", "  ")
	return os.WriteFile(indexPath, idxData, 0o644)
}
