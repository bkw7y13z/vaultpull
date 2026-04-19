package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Comment struct {
	Checksum  string    `json:"checksum"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type commentStore struct {
	Comments []Comment `json:"comments"`
}

func commentPath(snapshotPath string) string {
	dir := filepath.Dir(snapshotPath)
	return filepath.Join(dir, "comments.json")
}

func loadCommentStore(path string) (commentStore, error) {
	var store commentStore
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return store, err
	}
	return store, json.Unmarshal(data, &store)
}

func saveCommentStore(path string, store commentStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddComment(snapshotPath, checksum, text, author string) error {
	if snapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return fmt.Errorf("checksum is required")
	}
	if text == "" {
		return fmt.Errorf("comment text is required")
	}
	if author == "" {
		return fmt.Errorf("author is required")
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

	cp := commentPath(snapshotPath)
	store, err := loadCommentStore(cp)
	if err != nil {
		return err
	}
	store.Comments = append(store.Comments, Comment{
		Checksum:  checksum,
		Text:      text,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	})
	return saveCommentStore(cp, store)
}

func GetComments(snapshotPath, checksum string) ([]Comment, error) {
	if snapshotPath == "" {
		return nil, fmt.Errorf("snapshot path is required")
	}
	if checksum == "" {
		return nil, fmt.Errorf("checksum is required")
	}
	cp := commentPath(snapshotPath)
	store, err := loadCommentStore(cp)
	if err != nil {
		return nil, err
	}
	var result []Comment
	for _, c := range store.Comments {
		if c.Checksum == checksum {
			result = append(result, c)
		}
	}
	return result, nil
}
