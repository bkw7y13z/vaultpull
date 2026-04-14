package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp  time.Time `json:"timestamp"`
	Operation  string    `json:"operation"`
	VaultPath  string    `json:"vault_path"`
	OutputFile string    `json:"output_file"`
	KeysCount  int       `json:"keys_count"`
	Status     string    `json:"status"`
	Message    string    `json:"message,omitempty"`
}

// Logger writes structured audit entries to a file.
type Logger struct {
	path string
}

// NewLogger creates a Logger that appends to the given file path.
// Pass an empty string to disable logging (no-op logger).
func NewLogger(path string) *Logger {
	return &Logger{path: path}
}

// Log appends an audit entry to the log file.
// If the logger path is empty, this is a no-op.
func (l *Logger) Log(entry Entry) error {
	if l.path == "" {
		return nil
	}

	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("opening audit log: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshalling audit entry: %w", err)
	}

	_, err = fmt.Fprintf(f, "%s\n", data)
	if err != nil {
		return fmt.Errorf("writing audit entry: %w", err)
	}

	return nil
}

// Success is a convenience method to log a successful sync operation.
func (l *Logger) Success(vaultPath, outputFile string, keysCount int) error {
	return l.Log(Entry{
		Operation:  "sync",
		VaultPath:  vaultPath,
		OutputFile: outputFile,
		KeysCount:  keysCount,
		Status:     "success",
	})
}

// Failure is a convenience method to log a failed sync operation.
func (l *Logger) Failure(vaultPath, outputFile string, reason error) error {
	return l.Log(Entry{
		Operation:  "sync",
		VaultPath:  vaultPath,
		OutputFile: outputFile,
		Status:     "failure",
		Message:    reason.Error(),
	})
}
