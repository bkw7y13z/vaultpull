package snapshot_test

import (
	"testing"

	"github.com/nicholasgasior/vaultpull/internal/snapshot"
)

func TestComputeChecksum_Deterministic(t *testing.T) {
	secrets := map[string]string{
		"FOO": "bar",
		"BAZ": "qux",
	}
	c1 := snapshot.ComputeChecksum(secrets)
	c2 := snapshot.ComputeChecksum(secrets)
	if c1 != c2 {
		t.Errorf("checksum not deterministic: %q vs %q", c1, c2)
	}
}

func TestComputeChecksum_OrderIndependent(t *testing.T) {
	a := map[string]string{"KEY_A": "1", "KEY_B": "2"}
	b := map[string]string{"KEY_B": "2", "KEY_A": "1"}
	if snapshot.ComputeChecksum(a) != snapshot.ComputeChecksum(b) {
		t.Error("checksum should be order-independent")
	}
}

func TestComputeChecksum_DiffersOnValueChange(t *testing.T) {
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	if snapshot.ComputeChecksum(a) == snapshot.ComputeChecksum(b) {
		t.Error("checksum should differ when values change")
	}
}

func TestComputeChecksum_EmptyMap(t *testing.T) {
	c := snapshot.ComputeChecksum(map[string]string{})
	if c == "" {
		t.Error("expected non-empty checksum for empty map")
	}
}

func TestKeysFromSecrets_Sorted(t *testing.T) {
	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MIDDLE": "m"}
	keys := snapshot.KeysFromSecrets(secrets)
	expected := []string{"ALPHA", "MIDDLE", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
