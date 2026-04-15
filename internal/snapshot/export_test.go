package snapshot

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"
	"time"
)

func sampleEntries() []Entry {
	return []Entry{
		{
			Timestamp: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			Checksum:  "abc123",
			Keys:      []string{"DB_HOST", "API_KEY"},
		},
		{
			Timestamp: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
			Checksum:  "def456",
			Keys:      []string{"DB_HOST", "API_KEY", "SECRET"},
		},
	}
}

func TestExport_EmptyDestPath(t *testing.T) {
	err := Export(sampleEntries(), "", ExportOptions{Format: FormatText})
	if err == nil {
		t.Fatal("expected error for empty dest path")
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	tmp, _ := os.CreateTemp("", "export-*.txt")
	tmp.Close()
	defer os.Remove(tmp.Name())

	err := Export(sampleEntries(), tmp.Name(), ExportOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExport_TextFormat(t *testing.T) {
	tmp, _ := os.CreateTemp("", "export-*.txt")
	tmp.Close()
	defer os.Remove(tmp.Name())

	err := Export(sampleEntries(), tmp.Name(), ExportOptions{Format: FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	content := string(data)
	if !strings.Contains(content, "abc123") || !strings.Contains(content, "def456") {
		t.Errorf("expected checksums in output, got:\n%s", content)
	}
	if !strings.Contains(content, "API_KEY") {
		t.Errorf("expected keys in output, got:\n%s", content)
	}
}

func TestExport_CSVFormat(t *testing.T) {
	tmp, _ := os.CreateTemp("", "export-*.csv")
	tmp.Close()
	defer os.Remove(tmp.Name())

	err := Export(sampleEntries(), tmp.Name(), ExportOptions{Format: FormatCSV})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, _ := os.Open(tmp.Name())
	defer f.Close()
	records, _ := csv.NewReader(f).ReadAll()

	if len(records) != 3 { // header + 2 entries
		t.Fatalf("expected 3 CSV rows, got %d", len(records))
	}
	if records[0][0] != "timestamp" {
		t.Errorf("expected CSV header row")
	}
	if records[1][1] != "abc123" {
		t.Errorf("expected checksum abc123 in row 1")
	}
}

func TestExport_LimitEntries(t *testing.T) {
	tmp, _ := os.CreateTemp("", "export-*.txt")
	tmp.Close()
	defer os.Remove(tmp.Name())

	err := Export(sampleEntries(), tmp.Name(), ExportOptions{Format: FormatText, Limit: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	content := string(data)
	if strings.Contains(content, "abc123") {
		t.Error("expected only the last entry; abc123 should be excluded")
	}
	if !strings.Contains(content, "def456") {
		t.Error("expected def456 in limited export")
	}
}
