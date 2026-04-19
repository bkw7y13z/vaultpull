package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type WorkflowStep struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // pending, running, done, failed
	RunAt     time.Time `json:"run_at,omitempty"`
	Message   string    `json:"message,omitempty"`
}

type Workflow struct {
	Checksum  string         `json:"checksum"`
	CreatedBy string         `json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	Steps     []WorkflowStep `json:"steps"`
}

type workflowStore struct {
	Workflows []Workflow `json:"workflows"`
}

func workflowPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "workflows.json")
}

func loadWorkflowStore(path string) (workflowStore, error) {
	var store workflowStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveWorkflowStore(path string, store workflowStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func CreateWorkflow(snapshotPath, checksum, createdBy string, stepNames []string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if createdBy == "" {
		return errors.New("created_by is required")
	}
	if len(stepNames) == 0 {
		return errors.New("at least one step is required")
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

	steps := make([]WorkflowStep, len(stepNames))
	for i, name := range stepNames {
		steps[i] = WorkflowStep{Name: name, Status: "pending"}
	}

	wp := workflowPath(snapshotPath)
	store, err := loadWorkflowStore(wp)
	if err != nil {
		return err
	}
	store.Workflows = append(store.Workflows, Workflow{
		Checksum:  checksum,
		CreatedBy: createdBy,
		CreatedAt: time.Now().UTC(),
		Steps:     steps,
	})
	return saveWorkflowStore(wp, store)
}

func GetWorkflows(snapshotPath, checksum string) ([]Workflow, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}
	wp := workflowPath(snapshotPath)
	store, err := loadWorkflowStore(wp)
	if err != nil {
		return nil, err
	}
	var result []Workflow
	for _, w := range store.Workflows {
		if w.Checksum == checksum {
			result = append(result, w)
		}
	}
	return result, nil
}
