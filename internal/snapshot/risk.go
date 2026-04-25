package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

var validRiskLevels = map[RiskLevel]bool{
	RiskLow: true, RiskMedium: true, RiskHigh: true, RiskCritical: true,
}

type RiskEntry struct {
	Checksum  string    `json:"checksum"`
	Level     RiskLevel `json:"level"`
	Reason    string    `json:"reason"`
	AssessedBy string   `json:"assessed_by"`
	AssessedAt time.Time `json:"assessed_at"`
}

type riskStore struct {
	Entries []RiskEntry `json:"entries"`
}

func riskPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "risk.json")
}

func loadRiskStore(path string) (riskStore, error) {
	var store riskStore
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveRiskStore(path string, store riskStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AssessRisk(snapshotPath, checksum, reason, assessedBy string, level RiskLevel) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return fmt.Errorf("checksum is required")
	}
	if reason == "" {
		return fmt.Errorf("reason is required")
	}
	if assessedBy == "" {
		return fmt.Errorf("assessed_by is required")
	}
	if !validRiskLevels[level] {
		return fmt.Errorf("invalid risk level %q: must be low, medium, high, or critical", level)
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	if _, found := snap.FindByChecksum(checksum); !found {
		return fmt.Errorf("checksum %q not found in snapshot", checksum)
	}

	p := riskPath(snapshotPath)
	store, err := loadRiskStore(p)
	if err != nil {
		return err
	}

	for i, e := range store.Entries {
		if e.Checksum == checksum {
			store.Entries[i] = RiskEntry{Checksum: checksum, Level: level, Reason: reason, AssessedBy: assessedBy, AssessedAt: time.Now().UTC()}
			return saveRiskStore(p, store)
		}
	}

	store.Entries = append(store.Entries, RiskEntry{
		Checksum: checksum, Level: level, Reason: reason, AssessedBy: assessedBy, AssessedAt: time.Now().UTC(),
	})
	return saveRiskStore(p, store)
}

func GetRisk(snapshotPath, checksum string) (RiskEntry, bool, error) {
	if snapshotPath == "" {
		return RiskEntry{}, false, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return RiskEntry{}, false, fmt.Errorf("checksum is required")
	}
	store, err := loadRiskStore(riskPath(snapshotPath))
	if err != nil {
		return RiskEntry{}, false, err
	}
	for _, e := range store.Entries {
		if e.Checksum == checksum {
			return e, true, nil
		}
	}
	return RiskEntry{}, false, nil
}
