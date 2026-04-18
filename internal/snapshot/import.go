package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// ImportOptions controls how secrets are imported into a snapshot.
type ImportOptions struct {
	SnapshotPath string
	SourcePath   string
	Format       string // "json" or "env"
	Tag          string
	Overwrite    bool
}

// Import reads secrets from a file and adds them as a new snapshot entry.
func Import(opts ImportOptions) error {
	if opts.SnapshotPath == "" {
		return fmt.Errorf("snapshot path is required")
	}
	if opts.SourcePath == "" {
		return fmt.Errorf("source path is required")
	}
	if opts.Format == "" {
		opts.Format = "env"
	}

	data, err := os.ReadFile(opts.SourcePath)
	if err != nil {
		return fmt.Errorf("read source: %w", err)
	}

	var secrets map[string]string
	switch opts.Format {
	case "json":
		if err := json.Unmarshal(data, &secrets); err != nil {
			return fmt.Errorf("parse json: %w", err)
		}
	case "env":
		secrets, err = parseEnvFile(string(data))
		if err != nil {
			return fmt.Errorf("parse env: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	if len(secrets) == 0 {
		return fmt.Errorf("no secrets found in source file")
	}

	checksum := ComputeChecksum(secrets)

	snaps, err := Load(opts.SnapshotPath)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}

	if !opts.Overwrite {
		for _, e := range snaps {
			if e.Checksum == checksum {
				return fmt.Errorf("identical snapshot already exists (checksum: %s)", checksum[:8])
			}
		}
	}

	entry := Entry{
		Checksum:  checksum,
		Keys:      KeysFromSecrets(secrets),
		CreatedAt: time.Now().UTC(),
		Tag:       opts.Tag,
	}

	return Save(opts.SnapshotPath, append(snaps, entry))
}

func parseEnvFile(content string) (map[string]string, error) {
	result := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		result[key] = val
	}
	return result, nil
}
