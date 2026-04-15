package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RotateOptions controls how snapshot rotation behaves.
type RotateOptions struct {
	// ArchiveDir is the directory where rotated snapshots are moved.
	ArchiveDir string
	// MaxAgeDays removes archived files older than this many days. 0 disables age-based removal.
	MaxAgeDays int
}

// Rotate moves the current snapshot file to an archive directory, stamping it
// with the current UTC time, and optionally removes archived files that exceed
// MaxAgeDays. The source file is left intact if archiving fails.
func Rotate(snapshotPath string, opts RotateOptions) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path must not be empty")
	}
	if opts.ArchiveDir == "" {
		return fmt.Errorf("archive dir must not be empty")
	}

	if _, err := os.Stat(snapshotPath); os.IsNotExist(err) {
		return fmt.Errorf("snapshot file does not exist: %s", snapshotPath)
	}

	if err := os.MkdirAll(opts.ArchiveDir, 0o755); err != nil {
		return fmt.Errorf("creating archive dir: %w", err)
	}

	base := filepath.Base(snapshotPath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	stamp := time.Now().UTC().Format("20060102T150405Z")
	archiveName := fmt.Sprintf("%s_%s%s", name, stamp, ext)
	destPath := filepath.Join(opts.ArchiveDir, archiveName)

	src, err := os.ReadFile(snapshotPath)
	if err != nil {
		return fmt.Errorf("reading snapshot: %w", err)
	}
	if err := os.WriteFile(destPath, src, 0o644); err != nil {
		return fmt.Errorf("writing archive: %w", err)
	}

	if opts.MaxAgeDays > 0 {
		cutoff := time.Now().UTC().AddDate(0, 0, -opts.MaxAgeDays)
		entries, err := os.ReadDir(opts.ArchiveDir)
		if err != nil {
			return fmt.Errorf("reading archive dir: %w", err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			info, err := e.Info()
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				_ = os.Remove(filepath.Join(opts.ArchiveDir, e.Name()))
			}
		}
	}

	return nil
}
