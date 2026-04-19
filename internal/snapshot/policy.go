package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Policy struct {
	Checksum  string `json:"checksum"`
	MaxAge    int    `json:"max_age_days"`
	MinKeys   int    `json:"min_keys"`
	RequireTag bool  `json:"require_tag"`
	CreatedBy string `json:"created_by"`
}

type PolicyStore struct {
	Policies []Policy `json:"policies"`
}

func policyPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "policies.json")
}

func loadPolicyStore(path string) (PolicyStore, error) {
	var store PolicyStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func savePolicyStore(path string, store PolicyStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func SetPolicy(snapshotPath, checksum string, p Policy) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if p.CreatedBy == "" {
		return errors.New("created_by is required")
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
		return fmt.Errorf("checksum %q not found", checksum)
	}
	p.Checksum = checksum
	pp := policyPath(snapshotPath)
	store, err := loadPolicyStore(pp)
	if err != nil {
		return err
	}
	for i, existing := range store.Policies {
		if existing.Checksum == checksum {
			store.Policies[i] = p
			return savePolicyStore(pp, store)
		}
	}
	store.Policies = append(store.Policies, p)
	return savePolicyStore(pp, store)
}

func GetPolicy(snapshotPath, checksum string) (Policy, bool, error) {
	if snapshotPath == "" {
		return Policy{}, false, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return Policy{}, false, errors.New("checksum is required")
	}
	pp := policyPath(snapshotPath)
	store, err := loadPolicyStore(pp)
	if err != nil {
		return Policy{}, false, err
	}
	for _, p := range store.Policies {
		if p.Checksum == checksum {
			return p, true, nil
		}
	}
	return Policy{}, false, nil
}
