package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// WatchResult holds the outcome of a single watch poll cycle.
type WatchResult struct {
	Timestamp time.Time
	Checksum  string
	Changed   bool
	Diff      *DiffResult
}

// WatchOptions configures the Watch function.
type WatchOptions struct {
	SnapshotPath string
	Interval     time.Duration
	MaxCycles    int // 0 = run indefinitely
	OnChange     func(WatchResult)
	OnNoChange   func(WatchResult)
}

// Watch polls the snapshot file at the given interval and calls callbacks
// when changes are detected between consecutive snapshots.
func Watch(secrets func() (map[string]string, error), opts WatchOptions) error {
	if opts.SnapshotPath == "" {
		return errors.New("watch: snapshot path must not be empty")
	}
	if opts.Interval <= 0 {
		return errors.New("watch: interval must be positive")
	}
	if secrets == nil {
		return errors.New("watch: secrets fetch function must not be nil")
	}

	var prevChecksum string
	cycles := 0

	for {
		if opts.MaxCycles > 0 && cycles >= opts.MaxCycles {
			break
		}

		current, err := secrets()
		if err != nil {
			return fmt.Errorf("watch: failed to fetch secrets: %w", err)
		}

		checksum := ComputeChecksum(current)
		result := WatchResult{
			Timestamp: time.Now().UTC(),
			Checksum:  checksum,
			Changed:   checksum != prevChecksum && prevChecksum != "",
		}

		if result.Changed {
			snap, loadErr := Load(opts.SnapshotPath)
			if loadErr == nil {
				latest := snap.Latest()
				if latest != nil {
					dr := Diff(latest.Keys, KeysFromSecrets(current))
					result.Diff = &dr
				}
			}
			if opts.OnChange != nil {
				opts.OnChange(result)
			}
		} else if prevChecksum != "" {
			if opts.OnNoChange != nil {
				opts.OnNoChange(result)
			}
		}

		prevChecksum = checksum
		cycles++

		if opts.MaxCycles == 0 || cycles < opts.MaxCycles {
			time.Sleep(opts.Interval)
		}
	}

	return nil
}
