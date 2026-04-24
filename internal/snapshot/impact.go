package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ImpactEntry records which keys were affected and how many snapshots referenced them.
type ImpactEntry struct {
	Checksum  string    `json:"checksum"`
	Key       string    `json:"key"`
	RefCount  int       `json:"ref_count"`
	LastSeen  time.Time `json:"last_seen"`
	ChangedAt time.Time `json:"changed_at"`
}

// ImpactReport summarises the blast radius of a secret key change.
type ImpactReport struct {
	Key       string        `json:"key"`
	Entries   []ImpactEntry `json:"entries"`
	TotalRefs int           `json:"total_refs"`
}

func impactPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "impact_store.json")
}

func loadImpactStore(path string) ([]ImpactEntry, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []ImpactEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []ImpactEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func saveImpactStore(path string, entries []ImpactEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Impact analyses a snapshot file and builds an impact report for the given key.
func Impact(snapshotPath, key string) (*ImpactReport, error) {
	if snapshotPath == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}
	if key == "" {
		return nil, fmt.Errorf("key is required")
	}

	entries, err := Load(snapshotPath)
	if err != nil {
		return nil, fmt.Errorf("load snapshot: %w", err)
	}

	report := &ImpactReport{Key: key}
	for _, e := range entries {
		for _, k := range e.Keys {
			if k == key {
				report.Entries = append(report.Entries, ImpactEntry{
					Checksum:  e.Checksum,
					Key:       key,
					RefCount:  1,
					LastSeen:  e.CreatedAt,
					ChangedAt: e.CreatedAt,
				})
				report.TotalRefs++
				break
			}
		}
	}
	return report, nil
}
