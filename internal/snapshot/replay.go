package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// ReplayEvent represents a single step in a snapshot replay.
type ReplayEvent struct {
	Index     int
	Checksum  string
	Tag       string
	CreatedAt time.Time
	Keys      []string
	Diff      *DiffResult
}

// ReplayOptions controls how replay behaves.
type ReplayOptions struct {
	From string // checksum or tag
	To   string // checksum or tag (empty = latest)
	Dry  bool
}

// Replay walks snapshot entries between two refs, emitting events for each step.
func Replay(path string, opts ReplayOptions, fn func(ReplayEvent) error) error {
	if path == "" {
		return errors.New("snapshot path is required")
	}
	if opts.From == "" {
		return errors.New("from ref is required")
	}
	if fn == nil {
		return errors.New("event handler is required")
	}

	store, err := Load(path)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	if len(store.Entries) == 0 {
		return errors.New("snapshot is empty")
	}

	resolve := func(ref string) int {
		for i, e := range store.Entries {
			if e.Checksum == ref || e.Tag == ref {
				return i
			}
		}
		return -1
	}

	start := resolve(opts.From)
	if start == -1 {
		return fmt.Errorf("from ref %q not found", opts.From)
	}

	end := len(store.Entries) - 1
	if opts.To != "" {
		end = resolve(opts.To)
		if end == -1 {
			return fmt.Errorf("to ref %q not found", opts.To)
		}
	}

	if start > end {
		return errors.New("from ref is after to ref")
	}

	for i := start; i <= end; i++ {
		e := store.Entries[i]
		ev := ReplayEvent{
			Index:     i,
			Checksum:  e.Checksum,
			Tag:       e.Tag,
			CreatedAt: e.CreatedAt,
			Keys:      KeysFromSecrets(e.Secrets),
		}
		if i > start {
			prev := store.Entries[i-1]
			d := Diff(prev.Secrets, e.Secrets)
			ev.Diff = &d
		}
		if !opts.Dry {
			if err := fn(ev); err != nil {
				return fmt.Errorf("replay step %d: %w", i, err)
			}
		}
	}
	return nil
}
