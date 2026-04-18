package snapshot

import (
	"errors"
	"time"
)

// BlameEntry records who last modified a secret key and when.
type BlameEntry struct {
	Key       string    `json:"key"`
	Checksum  string    `json:"checksum"`
	Tag       string    `json:"tag,omitempty"`
	Note      string    `json:"note,omitempty"`
	ChangedAt time.Time `json:"changed_at"`
}

// Blame returns the most recent BlameEntry for each key across all snapshot entries.
// It walks entries from oldest to newest so the last write wins.
func Blame(path string) ([]BlameEntry, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}

	store, err := Load(path)
	if err != nil {
		return nil, err
	}

	if len(store.Entries) == 0 {
		return []BlameEntry{}, nil
	}

	// key -> latest blame
	blameMap := make(map[string]BlameEntry)

	for _, entry := range store.Entries {
		for _, key := range entry.Keys {
			blameMap[key] = BlameEntry{
				Key:       key,
				Checksum:  entry.Checksum,
				Tag:       entry.Tag,
				Note:      entry.Note,
				ChangedAt: entry.CreatedAt,
			}
		}
	}

	keys := make([]string, 0, len(blameMap))
	for k := range blameMap {
		keys = append(keys, k)
	}
	sortStrings(keys)

	result := make([]BlameEntry, 0, len(keys))
	for _, k := range keys {
		result = append(result, blameMap[k])
	}
	return result, nil
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
