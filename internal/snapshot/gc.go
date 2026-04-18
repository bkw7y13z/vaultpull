package snapshot

import (
	"errors"
	"time"
)

// GCOptions controls garbage collection behavior.
type GCOptions struct {
	SnapshotPath string
	DryRun       bool
	MaxAge       time.Duration
	KeepPinned   bool
	KeepTagged   bool
}

// GCResult summarizes what was removed.
type GCResult struct {
	Removed []string
	Kept    []string
}

// GC removes snapshot entries that are unpinned, untagged, and older than MaxAge.
func GC(opts GCOptions) (*GCResult, error) {
	if opts.SnapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if opts.MaxAge <= 0 {
		return nil, errors.New("max age must be positive")
	}

	snap, err := Load(opts.SnapshotPath)
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().UTC().Add(-opts.MaxAge)
	result := &GCResult{}
	var kept []Entry

	for _, e := range snap.Entries {
		age := e.CreatedAt.Before(cutoff)
		pinned := opts.KeepPinned && e.Pinned
		tagged := opts.KeepTagged && e.Tag != ""

		if age && !pinned && !tagged {
			result.Removed = append(result.Removed, e.Checksum)
		} else {
			kept = append(kept, e)
			result.Kept = append(result.Kept, e.Checksum)
		}
	}

	if !opts.DryRun {
		snap.Entries = kept
		if err := Save(opts.SnapshotPath, snap); err != nil {
			return nil, err
		}
	}

	return result, nil
}
