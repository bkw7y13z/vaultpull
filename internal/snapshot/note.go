package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type NoteStore struct {
	Notes map[string]string `json:"notes"`
}

func notePath(snapshotPath string) string {
	return filepath.Join(filepath.Dir(snapshotPath), "notes.json")
}

func loadNoteStore(snapshotPath string) (*NoteStore, error) {
	p := notePath(snapshotPath)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &NoteStore{Notes: map[string]string{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var store NoteStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	if store.Notes == nil {
		store.Notes = map[string]string{}
	}
	return &store, nil
}

func saveNoteStore(snapshotPath string, store *NoteStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(notePath(snapshotPath), data, 0600)
}

func AddNote(snapshotPath, checksum, note string) error {
	if snapshotPath == "" {
		return errors.New("snapshot path is required")
	}
	if checksum == "" {
		return errors.New("checksum is required")
	}
	if note == "" {
		return errors.New("note is required")
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
	store, err := loadNoteStore(snapshotPath)
	if err != nil {
		return err
	}
	store.Notes[checksum] = note
	return saveNoteStore(snapshotPath, store)
}

func GetNote(snapshotPath, checksum string) (string, error) {
	if snapshotPath == "" {
		return "", errors.New("snapshot path is required")
	}
	if checksum == "" {
		return "", errors.New("checksum is required")
	}
	store, err := loadNoteStore(snapshotPath)
	if err != nil {
		return "", err
	}
	n, ok := store.Notes[checksum]
	if !ok {
		return "", errors.New("no note found for checksum")
	}
	return n, nil
}
