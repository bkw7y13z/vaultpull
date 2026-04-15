package snapshot

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatch_EmptySnapshotPath(t *testing.T) {
	err := Watch(func() (map[string]string, error) { return nil, nil }, WatchOptions{
		Interval:  time.Millisecond,
		MaxCycles: 1,
	})
	if err == nil || err.Error() != "watch: snapshot path must not be empty" {
		t.Fatalf("expected empty path error, got %v", err)
	}
}

func TestWatch_ZeroInterval(t *testing.T) {
	err := Watch(func() (map[string]string, error) { return nil, nil }, WatchOptions{
		SnapshotPath: "/tmp/snap.json",
		MaxCycles:    1,
	})
	if err == nil {
		t.Fatal("expected interval error")
	}
}

func TestWatch_NilSecretsFn(t *testing.T) {
	err := Watch(nil, WatchOptions{
		SnapshotPath: "/tmp/snap.json",
		Interval:     time.Millisecond,
		MaxCycles:    1,
	})
	if err == nil {
		t.Fatal("expected nil secrets error")
	}
}

func TestWatch_FetchError(t *testing.T) {
	err := Watch(func() (map[string]string, error) {
		return nil, errors.New("vault down")
	}, WatchOptions{
		SnapshotPath: "/tmp/snap.json",
		Interval:     time.Millisecond,
		MaxCycles:    1,
	})
	if err == nil {
		t.Fatal("expected fetch error")
	}
}

func TestWatch_DetectsChange(t *testing.T) {
	var changeCalled atomic.Int32
	call := 0
	secrets := func() (map[string]string, error) {
		call++
		if call == 1 {
			return map[string]string{"KEY": "v1"}, nil
		}
		return map[string]string{"KEY": "v2"}, nil
	}

	err := Watch(secrets, WatchOptions{
		SnapshotPath: tmpPath(t),
		Interval:     time.Millisecond,
		MaxCycles:    2,
		OnChange: func(r WatchResult) {
			changeCalled.Add(1)
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changeCalled.Load() == 0 {
		t.Fatal("expected onChange to be called")
	}
}

func TestWatch_NoChangeCallback(t *testing.T) {
	var noChangeCalled atomic.Int32
	secrets := func() (map[string]string, error) {
		return map[string]string{"KEY": "stable"}, nil
	}

	err := Watch(secrets, WatchOptions{
		SnapshotPath: tmpPath(t),
		Interval:     time.Millisecond,
		MaxCycles:    3,
		OnNoChange: func(r WatchResult) {
			noChangeCalled.Add(1)
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if noChangeCalled.Load() == 0 {
		t.Fatal("expected onNoChange to be called")
	}
}
