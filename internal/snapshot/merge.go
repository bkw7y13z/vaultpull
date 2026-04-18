package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// MergeResult describes the outcome of merging two snapshots.
type MergeResult struct {
	Added   int
	Skipped int
	Total   int
}

// Merge combines entries from srcPath into dstPath.
// Entries already present in dst (by checksum) are skipped.
// Pinned entries in dst are never overwritten.
func Merge(dstPath, srcPath string) (MergeResult, error) {
	if dstPath == "" {
		return MergeResult{}, errors.New("destination snapshot path is required")
	}
	if srcPath == "" {
		return MergeResult{}, errors.New("source snapshot path is required")
	}

	dst, err := Load(dstPath)
	if err != nil {
		return MergeResult{}, fmt.Errorf("load destination: %w", err)
	}

	src, err := Load(srcPath)
	if err != nil {
		return MergeResult{}, fmt.Errorf("load source: %w", err)
	}

	existing := make(map[string]struct{}, len(dst.Entries))
	for _, e := range dst.Entries {
		existing[e.Checksum] = struct{}{}
	}

	var result MergeResult
	for _, e := range src.Entries {
		result.Total++
		if _, found := existing[e.Checksum]; found {
			result.Skipped++
			continue
		}
		e.CreatedAt = time.Now().UTC()
		dst.Entries = append(dst.Entries, e)
		existing[e.Checksum] = struct{}{}
		result.Added++
	}

	if err := Save(dstPath, dst); err != nil {
		return MergeResult{}, fmt.Errorf("save destination: %w", err)
	}

	return result, nil
}
