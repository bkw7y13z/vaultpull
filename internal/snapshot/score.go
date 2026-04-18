package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ScoreEntry holds a health score for a snapshot entry.
type ScoreEntry struct {
	Checksum  string    `json:"checksum"`
	Score     int       `json:"score"` // 0-100
	Reasons   []string  `json:"reasons"`
	ScoredAt  time.Time `json:"scored_at"`
}

func scorePath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "scores.json")
}

func loadScoreStore(path string) (map[string]ScoreEntry, error) {
	store := map[string]ScoreEntry{}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return nil, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveScoreStore(path string, store map[string]ScoreEntry) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Score computes and persists a health score for the given checksum.
func Score(snapshotPath, checksum string) (ScoreEntry, error) {
	if snapshotPath == "" {
		return ScoreEntry{}, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return ScoreEntry{}, fmt.Errorf("checksum is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return ScoreEntry{}, fmt.Errorf("load snapshot: %w", err)
	}

	var entry *Entry
	for i := range snap.Entries {
		if snap.Entries[i].Checksum == checksum {
			entry = &snap.Entries[i]
			break
		}
	}
	if entry == nil {
		return ScoreEntry{}, fmt.Errorf("checksum not found: %s", checksum)
	}

	score := 100
	var reasons []string

	if entry.Tag == "" {
		score -= 20
		reasons = append(reasons, "no tag assigned")
	}
	if len(entry.Keys) == 0 {
		score -= 30
		reasons = append(reasons, "no keys recorded")
	}
	if entry.CreatedAt.IsZero() {
		score -= 20
		reasons = append(reasons, "missing timestamp")
	}
	if time.Since(entry.CreatedAt) > 30*24*time.Hour {
		score -= 10
		reasons = append(reasons, "snapshot older than 30 days")
	}

	se := ScoreEntry{
		Checksum: checksum,
		Score:    score,
		Reasons:  reasons,
		ScoredAt: time.Now().UTC(),
	}

	sp := scorePath(snapshotPath)
	store, err := loadScoreStore(sp)
	if err != nil {
		return ScoreEntry{}, err
	}
	store[checksum] = se
	return se, saveScoreStore(sp, store)
}

// GetScore retrieves the stored score for a checksum.
func GetScore(snapshotPath, checksum string) (ScoreEntry, error) {
	if snapshotPath == "" {
		return ScoreEntry{}, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return ScoreEntry{}, fmt.Errorf("checksum is required")
	}
	store, err := loadScoreStore(scorePath(snapshotPath))
	if err != nil {
		return ScoreEntry{}, err
	}
	se, ok := store[checksum]
	if !ok {
		return ScoreEntry{}, fmt.Errorf("no score found for checksum: %s", checksum)
	}
	return se, nil
}
