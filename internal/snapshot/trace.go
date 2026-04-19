package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type TraceEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Checksum  string    `json:"checksum"`
	Operation string    `json:"operation"`
	Actor     string    `json:"actor"`
	Detail    string    `json:"detail,omitempty"`
}

type TraceStore struct {
	Events []TraceEvent `json:"events"`
}

func tracePath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "trace.json")
}

func loadTraceStore(path string) (TraceStore, error) {
	var store TraceStore
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveTraceStore(path string, store TraceStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddTrace(snapshotPath, checksum, operation, actor, detail string) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return fmt.Errorf("checksum is required")
	}
	if operation == "" {
		return fmt.Errorf("operation is required")
	}
	if actor == "" {
		return fmt.Errorf("actor is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	found := false
	for _, e := range snap {
		if e.Checksum == checksum {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("checksum not found: %s", checksum)
	}

	p := tracePath(snapshotPath)
	store, err := loadTraceStore(p)
	if err != nil {
		return err
	}
	store.Events = append(store.Events, TraceEvent{
		Timestamp: time.Now().UTC(),
		Checksum:  checksum,
		Operation: operation,
		Actor:     actor,
		Detail:    detail,
	})
	return saveTraceStore(p, store)
}

func GetTraces(snapshotPath, checksum string) ([]TraceEvent, error) {
	if snapshotPath == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return nil, fmt.Errorf("checksum is required")
	}
	p := tracePath(snapshotPath)
	store, err := loadTraceStore(p)
	if err != nil {
		return nil, err
	}
	var result []TraceEvent
	for _, e := range store.Events {
		if e.Checksum == checksum {
			result = append(result, e)
		}
	}
	return result, nil
}
