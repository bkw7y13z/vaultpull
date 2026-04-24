package snapshot

import (
	"errors"
	"fmt"
	"strings"
)

// ResolveResult holds the resolved checksum and the method used to find it.
type ResolveResult struct {
	Checksum string
	Method   string // "checksum", "tag", "bookmark", "label"
}

// Resolve attempts to find a snapshot entry by ref, which may be a full
// checksum, a tag name, a bookmark label, or a label name. Returns the
// resolved checksum and the method used, or an error if not found.
func Resolve(snapshotPath, ref string) (ResolveResult, error) {
	if snapshotPath == "" {
		return ResolveResult{}, errors.New("snapshot path is required")
	}
	if ref == "" {
		return ResolveResult{}, errors.New("ref is required")
	}

	// Try exact checksum match first.
	entries, err := Load(snapshotPath)
	if err != nil {
		return ResolveResult{}, fmt.Errorf("load snapshot: %w", err)
	}
	for _, e := range entries {
		if strings.EqualFold(e.Checksum, ref) {
			return ResolveResult{Checksum: e.Checksum, Method: "checksum"}, nil
		}
	}

	// Try tag lookup.
	tagEntry, err := FindByTag(snapshotPath, ref)
	if err == nil && tagEntry != nil {
		return ResolveResult{Checksum: tagEntry.Checksum, Method: "tag"}, nil
	}

	// Try bookmark lookup.
	bm, err := GetBookmark(snapshotPath, ref)
	if err == nil && bm != "" {
		return ResolveResult{Checksum: bm, Method: "bookmark"}, nil
	}

	// Try label lookup.
	lbl, err := GetLabel(snapshotPath, ref)
	if err == nil && lbl != "" {
		return ResolveResult{Checksum: lbl, Method: "label"}
	}

	return ResolveResult{}, fmt.Errorf("ref %q not found as checksum, tag, bookmark, or label", ref)
}
