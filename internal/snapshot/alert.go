package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type AlertSeverity string

const (
	AlertInfo     AlertSeverity = "info"
	AlertWarning  AlertSeverity = "warning"
	AlertCritical AlertSeverity = "critical"
)

type Alert struct {
	Checksum  string        `json:"checksum"`
	Severity  AlertSeverity `json:"severity"`
	Message   string        `json:"message"`
	CreatedBy string        `json:"created_by"`
	CreatedAt time.Time     `json:"created_at"`
}

type alertStore struct {
	Alerts []Alert `json:"alerts"`
}

func alertPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "alerts.json")
}

func loadAlertStore(path string) (alertStore, error) {
	var store alertStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveAlertStore(path string, store alertStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddAlert(snapshotPath, checksum, message, createdBy string, severity AlertSeverity) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if message == "" {
		return errors.New("message is required")
	}
	if createdBy == "" {
		return errors.New("created_by is required")
	}
	if severity != AlertInfo && severity != AlertWarning && severity != AlertCritical {
		return errors.New("invalid severity: must be info, warning, or critical")
	}
	snap, err := Load(snapshotPath)
	if err != nil {
		return err
	}
	found := false
	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			found = true
			break
		}
	}
	if !found {
		return errors.New("checksum not found in snapshot")
	}
	ap := alertPath(snapshotPath)
	store, err := loadAlertStore(ap)
	if err != nil {
		return err
	}
	store.Alerts = append(store.Alerts, Alert{
		Checksum:  checksum,
		Severity:  severity,
		Message:   message,
		CreatedBy: createdBy,
		CreatedAt: time.Now().UTC(),
	})
	return saveAlertStore(ap, store)
}

func GetAlerts(snapshotPath, checksum string) ([]Alert, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	ap := alertPath(snapshotPath)
	store, err := loadAlertStore(ap)
	if err != nil {
		return nil, err
	}
	var result []Alert
	for _, a := range store.Alerts {
		if a.Checksum == checksum {
			result = append(result, a)
		}
	}
	return result, nil
}
