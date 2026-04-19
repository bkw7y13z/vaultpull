package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAlert_Integration_AddAndRetrieve(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	entry := Entry{
		Checksum:  "integ001",
		Keys:      []string{"DB_PASS", "API_KEY"},
		CreatedAt: time.Now().UTC(),
	}
	data, _ := json.Marshal(Snapshot{Entries: []Entry{entry}})
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	if err := AddAlert(path, "integ001", "value drift detected", "monitor", AlertWarning); err != nil {
		t.Fatalf("AddAlert failed: %v", err)
	}
	if err := AddAlert(path, "integ001", "secret expired", "scheduler", AlertCritical); err != nil {
		t.Fatalf("AddAlert second failed: %v", err)
	}

	alerts, err := GetAlerts(path, "integ001")
	if err != nil {
		t.Fatalf("GetAlerts failed: %v", err)
	}
	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}

	severities := map[AlertSeverity]bool{}
	for _, a := range alerts {
		severities[a.Severity] = true
	}
	if !severities[AlertWarning] || !severities[AlertCritical] {
		t.Error("expected both warning and critical alerts")
	}

	none, err := GetAlerts(path, "noexist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(none) != 0 {
		t.Errorf("expected 0 alerts for unknown checksum, got %d", len(none))
	}
}
