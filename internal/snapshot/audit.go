package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Checksum  string    `json:"checksum,omitempty"`
	Actor     string    `json:"actor,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

func auditPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "audit.json")
}

func loadAuditLog(path string) (*AuditLog, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, err
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, err
	}
	return &log, nil
}

func saveAuditLog(path string, log *AuditLog) error {
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func RecordAudit(snapshotPath, action, checksum, actor, detail string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if action == "" {
		return errors.New("action is required")
	}
	p := auditPath(snapshotPath)
	log, err := loadAuditLog(p)
	if err != nil {
		return err
	}
	log.Events = append(log.Events, AuditEvent{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Checksum:  checksum,
		Actor:     actor,
		Detail:    detail,
	})
	return saveAuditLog(p, log)
}

func GetAuditLog(snapshotPath string) ([]AuditEvent, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	log, err := loadAuditLog(auditPath(snapshotPath))
	if err != nil {
		return nil, err
	}
	return log.Events, nil
}
