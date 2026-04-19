package snapshot

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeWorkflowSnapshot(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "snap-*.json")
	if err != nil {
		t.Fatal(err)
	}
	snap := Snapshot{Entries: entries}
	if err := json.NewEncoder(f).Encode(snap); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestCreateWorkflow_EmptyPath(t *testing.T) {
	err := CreateWorkflow("", "abc", "ci", []string{"build"})
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestCreateWorkflow_EmptyChecksum(t *testing.T) {
	path := writeWorkflowSnapshot(t, nil)
	err := CreateWorkflow(path, "", "ci", []string{"build"})
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestCreateWorkflow_EmptyCreatedBy(t *testing.T) {
	path := writeWorkflowSnapshot(t, nil)
	err := CreateWorkflow(path, "abc", "", []string{"build"})
	if err == nil || err.Error() != "created_by is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestCreateWorkflow_NoSteps(t *testing.T) {
	path := writeWorkflowSnapshot(t, nil)
	err := CreateWorkflow(path, "abc", "ci", nil)
	if err == nil || err.Error() != "at least one step is required" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestCreateWorkflow_ChecksumNotFound(t *testing.T) {
	path := writeWorkflowSnapshot(t, []Entry{})
	err := CreateWorkflow(path, "missing", "ci", []string{"build"})
	if err == nil || err.Error() != "checksum not found in snapshot" {
		t.Fatalf("expected error, got %v", err)
	}
}

func TestCreateWorkflow_Success(t *testing.T) {
	entry := Entry{Checksum: "abc123", Keys: []string{"K"}, At: time.Now().UTC()}
	path := writeWorkflowSnapshot(t, []Entry{entry})

	err := CreateWorkflow(path, "abc123", "ci-bot", []string{"lint", "deploy"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	workflows, err := GetWorkflows(path, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(workflows) != 1 {
		t.Fatalf("expected 1 workflow, got %d", len(workflows))
	}
	w := workflows[0]
	if w.CreatedBy != "ci-bot" {
		t.Errorf("expected ci-bot, got %s", w.CreatedBy)
	}
	if len(w.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(w.Steps))
	}
	for _, s := range w.Steps {
		if s.Status != "pending" {
			t.Errorf("expected pending, got %s", s.Status)
		}
	}
}

func TestGetWorkflows_EmptyPath(t *testing.T) {
	_, err := GetWorkflows("", "abc")
	if err == nil {
		t.Fatal("expected error")
	}
}
