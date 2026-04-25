package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type SignatureRecord struct {
	Checksum  string    `json:"checksum"`
	SignedBy  string    `json:"signed_by"`
	SignedAt  time.Time `json:"signed_at"`
	PublicKey string    `json:"public_key"`
	Comment   string    `json:"comment,omitempty"`
}

type signatureStore struct {
	Signatures []SignatureRecord `json:"signatures"`
}

func signaturePath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "signatures.json")
}

func loadSignatureStore(path string) (signatureStore, error) {
	var store signatureStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveSignatureStore(path string, store signatureStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// AddSignature records a signature for a snapshot entry identified by checksum.
func AddSignature(snapshotPath, checksum, signedBy, publicKey, comment string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if signedBy == "" {
		return errors.New("signed_by is required")
	}
	if publicKey == "" {
		return errors.New("public_key is required")
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
		return fmt.Errorf("checksum %q not found in snapshot", checksum)
	}

	sp := signaturePath(snapshotPath)
	store, err := loadSignatureStore(sp)
	if err != nil {
		return fmt.Errorf("load signature store: %w", err)
	}

	store.Signatures = append(store.Signatures, SignatureRecord{
		Checksum:  checksum,
		SignedBy:  signedBy,
		SignedAt:  time.Now().UTC(),
		PublicKey: publicKey,
		Comment:   comment,
	})

	return saveSignatureStore(sp, store)
}

// GetSignatures returns all signature records for a given checksum.
func GetSignatures(snapshotPath, checksum string) ([]SignatureRecord, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return nil, errors.New("checksum is required")
	}

	sp := signaturePath(snapshotPath)
	store, err := loadSignatureStore(sp)
	if err != nil {
		return nil, fmt.Errorf("load signature store: %w", err)
	}

	var results []SignatureRecord
	for _, r := range store.Signatures {
		if r.Checksum == checksum {
			results = append(results, r)
		}
	}
	return results, nil
}
