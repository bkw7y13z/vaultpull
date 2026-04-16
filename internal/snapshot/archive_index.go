package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LoadArchiveIndex reads the archive index from the given directory.
func LoadArchiveIndex(archiveDir string) (*ArchiveIndex, error) {
	if archiveDir == "" {
		return nil, errors.New("archive dir is required")
	}
	idxPath := filepath.Join(archiveDir, "archive_index.json")
	data, err := os.ReadFile(idxPath)
	if os.IsNotExist(err) {
		return &ArchiveIndex{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}
	var idx ArchiveIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse index: %w", err)
	}
	return &idx, nil
}

// FindInArchive looks up an archived entry by checksum.
func FindInArchive(archiveDir, checksum string) (*Entry, error) {
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	idx, err := LoadArchiveIndex(archiveDir)
	if err != nil {
		return nil, err
	}
	for _, ref := range idx.Entries {
		if ref.Checksum == checksum {
			data, err := os.ReadFile(ref.FilePath)
			if err != nil {
				return nil, fmt.Errorf("read archived entry: %w", err)
			}
			var e Entry
			if err := json.Unmarshal(data, &e); err != nil {
				return nil, fmt.Errorf("parse archived entry: %w", err)
			}
			return &e, nil
		}
	}
	return nil, fmt.Errorf("checksum %s not found in archive", checksum)
}
