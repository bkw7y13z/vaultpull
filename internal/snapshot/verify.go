package snapshot

import (
	"errors"
	"fmt"
	"time"
)

// VerifyResult holds the outcome of verifying a single snapshot entry.
type VerifyResult struct {
	Checksum  string
	Tag       string
	Timestamp time.Time
	Valid     bool
	Reason    string
}

// Verify checks all entries in the snapshot file for integrity.
// It validates checksums are non-empty, timestamps are set, and no duplicate checksums exist.
func Verify(snapshotPath string) ([]VerifyResult, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot: %w", err)
	}

	if len(entries) == 0 {
		return []VerifyResult{}, nil
	}

	seen := make(map[string]bool)
	results := make([]VerifyResult, 0, len(entries))

	for _, e := range entries {
		r := VerifyResult{
			Checksum:  e.Checksum,
			Tag:       e.Tag,
			Timestamp: e.CreatedAt,
			Valid:     true,
		}

		switch {
		case e.Checksum == "":
			r.Valid = false
			r.Reason = "missing checksum"
		case e.CreatedAt.IsZero():
			r.Valid = false
			r.Reason = "missing timestamp"
		case seen[e.Checksum]:
			r.Valid = false
			r.Reason = "duplicate checksum"
		default:
			seen[e.Checksum] = true
		}

		results = append(results, r)
	}

	return results, nil
}
