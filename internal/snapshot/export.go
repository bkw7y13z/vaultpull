package snapshot

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"time"
)

// ExportFormat defines the output format for snapshot exports.
type ExportFormat string

const (
	FormatCSV  ExportFormat = "csv"
	FormatText ExportFormat = "text"
)

// ExportOptions configures how a snapshot history is exported.
type ExportOptions struct {
	Format ExportFormat
	Limit  int
}

// Export writes snapshot history to the given file path in the specified format.
func Export(snapshots []Entry, destPath string, opts ExportOptions) error {
	if destPath == "" {
		return fmt.Errorf("export destination path must not be empty")
	}

	entries := snapshots
	if opts.Limit > 0 && opts.Limit < len(entries) {
		entries = entries[len(entries)-opts.Limit:]
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer f.Close()

	switch opts.Format {
	case FormatCSV:
		return exportCSV(f, entries)
	case FormatText, "":
		return exportText(f, entries)
	default:
		return fmt.Errorf("unsupported export format: %s", opts.Format)
	}
}

func exportCSV(f *os.File, entries []Entry) error {
	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"timestamp", "checksum", "keys"}); err != nil {
		return err
	}

	for _, e := range entries {
		keys := sortedKeys(e.Keys)
		row := []string{
			e.Timestamp.Format(time.RFC3339),
			e.Checksum,
			formatKeys(keys),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	return w.Error()
}

func exportText(f *os.File, entries []Entry) error {
	for _, e := range entries {
		keys := sortedKeys(e.Keys)
		line := fmt.Sprintf("[%s] checksum=%s keys=[%s]\n",
			e.Timestamp.Format(time.RFC3339),
			e.Checksum,
			formatKeys(keys),
		)
		if _, err := fmt.Fprint(f, line); err != nil {
			return err
		}
	}
	return nil
}

func sortedKeys(keys []string) []string {
	copy_ := append([]string{}, keys...)
	sort.Strings(copy_)
	return copy_
}

func formatKeys(keys []string) string {
	result := ""
	for i, k := range keys {
		if i > 0 {
			result += ","
		}
		result += k
	}
	return result
}
