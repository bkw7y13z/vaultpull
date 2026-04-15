package snapshot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiff_NoChanges(t *testing.T) {
	prev := map[string]string{"KEY_A": "val1", "KEY_B": "val2"}
	curr := map[string]string{"KEY_A": "val1", "KEY_B": "val2"}

	result := Diff(prev, curr)

	assert.False(t, result.HasChanges())
	assert.Empty(t, result.Added)
	assert.Empty(t, result.Removed)
	assert.Empty(t, result.Changed)
}

func TestDiff_AddedKeys(t *testing.T) {
	prev := map[string]string{"KEY_A": "val1"}
	curr := map[string]string{"KEY_A": "val1", "KEY_B": "val2", "KEY_C": "val3"}

	result := Diff(prev, curr)

	assert.True(t, result.HasChanges())
	assert.Equal(t, []string{"KEY_B", "KEY_C"}, result.Added)
	assert.Empty(t, result.Removed)
	assert.Empty(t, result.Changed)
}

func TestDiff_RemovedKeys(t *testing.T) {
	prev := map[string]string{"KEY_A": "val1", "KEY_B": "val2"}
	curr := map[string]string{"KEY_A": "val1"}

	result := Diff(prev, curr)

	assert.True(t, result.HasChanges())
	assert.Empty(t, result.Added)
	assert.Equal(t, []string{"KEY_B"}, result.Removed)
	assert.Empty(t, result.Changed)
}

func TestDiff_ChangedKeys(t *testing.T) {
	prev := map[string]string{"KEY_A": "old", "KEY_B": "same"}
	curr := map[string]string{"KEY_A": "new", "KEY_B": "same"}

	result := Diff(prev, curr)

	assert.True(t, result.HasChanges())
	assert.Empty(t, result.Added)
	assert.Empty(t, result.Removed)
	assert.Equal(t, []string{"KEY_A"}, result.Changed)
}

func TestDiff_MixedChanges(t *testing.T) {
	prev := map[string]string{"KEY_A": "old", "KEY_B": "gone"}
	curr := map[string]string{"KEY_A": "new", "KEY_C": "fresh"}

	result := Diff(prev, curr)

	assert.True(t, result.HasChanges())
	assert.Equal(t, []string{"KEY_C"}, result.Added)
	assert.Equal(t, []string{"KEY_B"}, result.Removed)
	assert.Equal(t, []string{"KEY_A"}, result.Changed)
}

func TestDiff_EmptyBoth(t *testing.T) {
	result := Diff(map[string]string{}, map[string]string{})

	assert.False(t, result.HasChanges())
}

func TestDiff_SortedOutput(t *testing.T) {
	prev := map[string]string{}
	curr := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}

	result := Diff(prev, curr)

	assert.Equal(t, []string{"APPLE", "MANGO", "ZEBRA"}, result.Added)
}
