package env

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteEnvFile_EmptyPath(t *testing.T) {
	err := WriteEnvFile(map[string]string{"KEY": "val"}, WriteOptions{})
	if err == nil || !strings.Contains(err.Error(), "output path") {
		t.Fatalf("expected output path error, got %v", err)
	}
}

func TestWriteEnvFile_NoOverwriteExisting(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(tmp, []byte("EXISTING=1\n"), 0600)

	err := WriteEnvFile(map[string]string{"KEY": "val"}, WriteOptions{OutputPath: tmp})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected already-exists error, got %v", err)
	}
}

func TestWriteEnvFile_Overwrite(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(tmp, []byte("OLD=1\n"), 0600)

	err := WriteEnvFile(map[string]string{"NEW": "value"}, WriteOptions{OutputPath: tmp, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "NEW=value") {
		t.Errorf("expected NEW=value in output, got: %s", data)
	}
	if strings.Contains(string(data), "OLD") {
		t.Errorf("expected OLD to be removed after overwrite")
	}
}

func TestWriteEnvFile_Append(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	_ = os.WriteFile(tmp, []byte("EXISTING=1\n"), 0600)

	err := WriteEnvFile(map[string]string{"ADDED": "2"}, WriteOptions{OutputPath: tmp, Append: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "EXISTING=1") || !strings.Contains(string(data), "ADDED=2") {
		t.Errorf("expected both keys, got: %s", data)
	}
}

func TestWriteEnvFile_QuotesSpecialValues(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")

	err := WriteEnvFile(map[string]string{"MSG": "hello world"}, WriteOptions{OutputPath: tmp, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), `MSG="hello world"`) {
		t.Errorf("expected quoted value, got: %s", data)
	}
}

func TestWriteEnvFile_SortedOutput(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), ".env")
	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a", "MIDDLE": "m"}

	err := WriteEnvFile(secrets, WriteOptions{OutputPath: tmp, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "ALPHA") || !strings.HasPrefix(lines[2], "ZEBRA") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
