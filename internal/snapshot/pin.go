package snapshot

import (
	"errors"
	"time"
)

// PinEntry records a pinned snapshot checksum with a reason.
type PinEntry struct {
	Checksum  string    `json:"checksum"`
	Reason    string    `json:"reason"`
	PinnedAt  time.Time `json:"pinned_at"`
}

// Pin marks a snapshot entry as pinned so prune/rotate will not remove it.
func Pin(path, checksum, reason string) error {
	if path == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if reason == "" {
		return errors.New("reason is required")
	}

	snap, err := Load(path)
	if err != nil {
		return err
	}

	found := false
	for i, e := range snap.Entries {
		if e.Checksum == checksum {
			if snap.Entries[i].Annotations == nil {
				snap.Entries[i].Annotations = map[string]string{}
			}
			snap.Entries[i].Annotations["pinned"] = "true"
			snap.Entries[i].Annotations["pin_reason"] = reason
			found = true
			break
		}
	}
	if !found {
		return errors.New("checksum not found in snapshot")
	}

	return Save(path, snap)
}

// Unpin removes the pinned marker from a snapshot entry.
func Unpin(path, checksum string) error {
	if path == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}

	snap, err := Load(path)
	if err != nil {
		return err
	}

	found := false
	for i, e := range snap.Entries {
		if e.Checksum == checksum {
			if snap.Entries[i].Annotations != nil {
				delete(snap.Entries[i].Annotations, "pinned")
				delete(snap.Entries[i].Annotations, "pin_reason")
			}
			found = true
			break
		}
	}
	if !found {
		return errors.New("checksum not found in snapshot")
	}

	return Save(path, snap)
}

// IsPinned returns true if the entry with the given checksum is pinned.
func IsPinned(path, checksum string) (bool, error) {
	snap, err := Load(path)
	if err != nil {
		return false, err
	}
	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			return e.Annotations["pinned"] == "true", nil
		}
	}
	return false, errors.New("checksum not found in snapshot")
}
