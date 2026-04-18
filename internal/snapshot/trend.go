package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// TrendPoint represents the secret count at a point in time.
type TrendPoint struct {
	At       time.Time
	Checksum string
	Tag      string
	KeyCount int
}

// Trend returns a time-ordered series of key counts across snapshot entries.
func Trend(path string, since time.Time) ([]TrendPoint, error) {
	if path == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}

	store, err := Load(path)
	if err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}

	if len(store.Entries) == 0 {
		return []TrendPoint{}, nil
	}

	var points []TrendPoint
	for _, e := range store.Entries {
		if !since.IsZero() && e.CreatedAt.Before(since) {
			continue
		}
		points = append(points, TrendPoint{
			At:       e.CreatedAt,
			Checksum: e.Checksum,
			Tag:      e.Tag,
			KeyCount: len(e.Keys),
		})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].At.Before(points[j].At)
	})

	return points, nil
}
