package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Bookmark struct {
	Label    string `json:"label"`
	Checksum string `json:"checksum"`
	Note     string `json:"note,omitempty"`
}

type BookmarkStore struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}

func bookmarkPath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "bookmarks.json")
}

func loadBookmarkStore(path string) (BookmarkStore, error) {
	var store BookmarkStore
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveBookmarkStore(path string, store BookmarkStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddBookmark(snapshotPath, label, checksum, note string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if label == "" {
		return errors.New("label is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
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
	bp := bookmarkPath(snapshotPath)
	store, err := loadBookmarkStore(bp)
	if err != nil {
		return err
	}
	for _, b := range store.Bookmarks {
		if b.Label == label {
			return errors.New("bookmark label already exists")
		}
	}
	store.Bookmarks = append(store.Bookmarks, Bookmark{Label: label, Checksum: checksum, Note: note})
	return saveBookmarkStore(bp, store)
}

func GetBookmark(snapshotPath, label string) (Bookmark, error) {
	if snapshotPath == "" {
		return Bookmark{}, errors.New("snapshot path is required")
	}
	if label == "" {
		return Bookmark{}, errors.New("label is required")
	}
	bp := bookmarkPath(snapshotPath)
	store, err := loadBookmarkStore(bp)
	if err != nil {
		return Bookmark{}, err
	}
	for _, b := range store.Bookmarks {
		if b.Label == label {
			return b, nil
		}
	}
	return Bookmark{}, errors.New("bookmark not found")
}

func ListBookmarks(snapshotPath string) ([]Bookmark, error) {
	if snapshotPath == "" {
		return nil, errors.New("snapshot path is required")
	}
	bp := bookmarkPath(snapshotPath)
	store, err := loadBookmarkStore(bp)
	if err != nil {
		return nil, err
	}
	return store.Bookmarks, nil
}
