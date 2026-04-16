package snapshot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// SignEntry attaches an HMAC-SHA256 signature to a snapshot entry identified by checksum.
func SignEntry(path, checksum, secret string) error {
	if path == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if secret == "" {
		return errors.New("signing secret is required")
	}

	snap, err := Load(path)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	idx := -1
	for i, e := range snap.Entries {
		if e.Checksum == checksum {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("entry with checksum %q not found", checksum)
	}

	sig := computeHMAC(checksum, secret)
	if snap.Entries[idx].Metadata == nil {
		snap.Entries[idx].Metadata = map[string]string{}
	}
	snap.Entries[idx].Metadata["signature"] = sig

	return Save(path, snap)
}

// VerifySignature checks the HMAC-SHA256 signature of an entry.
func VerifySignature(path, checksum, secret string) (bool, error) {
	if path == "" {
		return false, errors.New("snapshot path is required")
	}
	if checksum == "" {
		return false, errors.New("checksum is required")
	}
	if secret == "" {
		return false, errors.New("signing secret is required")
	}

	snap, err := Load(path)
	if err != nil {
		return false, fmt.Errorf("load snapshot: %w", err)
	}

	for _, e := range snap.Entries {
		if e.Checksum == checksum {
			if e.Metadata == nil {
				return false, nil
			}
			stored, ok := e.Metadata["signature"]
			if !ok {
				return false, nil
			}
			expected := computeHMAC(checksum, secret)
			return hmac.Equal([]byte(stored), []byte(expected)), nil
		}
	}
	return false, fmt.Errorf("entry with checksum %q not found", checksum)
}

func computeHMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
