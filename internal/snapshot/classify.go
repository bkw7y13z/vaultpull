package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Classification struct {
	Checksum  string    `json:"checksum"`
	Category  string    `json:"category"`
	Sensitive bool      `json:"sensitive"`
	SetBy     string    `json:"set_by"`
	SetAt     time.Time `json:"set_at"`
}

type classifyStore struct {
	Entries map[string]Classification `json:"entries"`
}

func classifyPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "classify.json")
}

func loadClassifyStore(path string) (classifyStore, error) {
	var store classifyStore
	store.Entries = make(map[string]Classification)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveClassifyStore(path string, store classifyStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Classify assigns a category and sensitivity flag to a snapshot entry.
func Classify(snapshotPath, checksum, category, setBy string, sensitive bool) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return fmt.Errorf("checksum is required")
	}
	if category == "" {
		return fmt.Errorf("category is required")
	}
	if setBy == "" {
		return fmt.Errorf("set_by is required")
	}

	snap, err := Load(snapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	found := false
	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("checksum not found: %s", checksum)
	}

	store, err := loadClassifyStore(classifyPath(snapshotPath))
	if err != nil {
		return err
	}
	store.Entries[checksum] = Classification{
		Checksum:  checksum,
		Category:  category,
		Sensitive: sensitive,
		SetBy:     setBy,
		SetAt:     time.Now().UTC(),
	}
	return saveClassifyStore(classifyPath(snapshotPath), store)
}

// GetClassification retrieves the classification for a given checksum.
func GetClassification(snapshotPath, checksum string) (Classification, bool, error) {
	if snapshotPath == "" {
		return Classification{}, false, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return Classification{}, false, fmt.Errorf("checksum is required")
	}
	store, err := loadClassifyStore(classifyPath(snapshotPath))
	if err != nil {
		return Classification{}, false, err
	}
	c, ok := store.Entries[checksum]
	return c, ok, nil
}
