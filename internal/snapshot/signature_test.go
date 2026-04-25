package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeSignatureSnapshot(t *testing.T, dir string) string {
	t.Helper()
	snap := Snapshot{
		Entries: []Entry{
			{Checksum: "abc123", Keys: []string{"KEY"}, CreatedAt: time.Now().UTC()},
		},
	}
	data, _ := json.Marshal(snap)
	p := filepath.Join(dir, "snapshot.json")
	_ = os.WriteFile(p, data, 0600)
	return p
}

func TestAddSignature_EmptyPath(t *testing.T) {
	err := AddSignature("", "abc123", "alice", "pubkey", "")
	if err == nil || err.Error() != "snapshot path is required" {
		t.Fatalf("expected snapshot path error, got %v", err)
	}
}

func TestAddSignature_EmptyChecksum(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)
	err := AddSignature(p, "", "alice", "pubkey", "")
	if err == nil || err.Error() != "checksum is required" {
		t.Fatalf("expected checksum error, got %v", err)
	}
}

func TestAddSignature_EmptySignedBy(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)
	err := AddSignature(p, "abc123", "", "pubkey", "")
	if err == nil || err.Error() != "signed_by is required" {
		t.Fatalf("expected signed_by error, got %v", err)
	}
}

func TestAddSignature_EmptyPublicKey(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)
	err := AddSignature(p, "abc123", "alice", "", "")
	if err == nil || err.Error() != "public_key is required" {
		t.Fatalf("expected public_key error, got %v", err)
	}
}

func TestAddSignature_ChecksumNotFound(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)
	err := AddSignature(p, "notexist", "alice", "pubkey", "")
	if err == nil {
		t.Fatal("expected error for missing checksum")
	}
}

func TestAddSignature_Success(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)

	err := AddSignature(p, "abc123", "alice", "ssh-ed25519 AAAA...", "release sign-off")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sigs, err := GetSignatures(p, "abc123")
	if err != nil {
		t.Fatalf("get signatures: %v", err)
	}
	if len(sigs) != 1 {
		t.Fatalf("expected 1 signature, got %d", len(sigs))
	}
	if sigs[0].SignedBy != "alice" {
		t.Errorf("expected signed_by=alice, got %s", sigs[0].SignedBy)
	}
	if sigs[0].Comment != "release sign-off" {
		t.Errorf("expected comment, got %s", sigs[0].Comment)
	}
}

func TestAddSignature_Appends(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)

	_ = AddSignature(p, "abc123", "alice", "key1", "first")
	_ = AddSignature(p, "abc123", "bob", "key2", "second")

	sigs, err := GetSignatures(p, "abc123")
	if err != nil {
		t.Fatalf("get signatures: %v", err)
	}
	if len(sigs) != 2 {
		t.Fatalf("expected 2 signatures, got %d", len(sigs))
	}
}

func TestGetSignatures_EmptySnapshot(t *testing.T) {
	dir := t.TempDir()
	p := writeSignatureSnapshot(t, dir)

	sigs, err := GetSignatures(p, "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sigs) != 0 {
		t.Errorf("expected 0 signatures, got %d", len(sigs))
	}
}
