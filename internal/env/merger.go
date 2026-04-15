package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// MergeResult holds the outcome of a merge operation.
type MergeResult struct {
	Added   []string
	Updated []string
	Skipped []string
}

// MergeEnvFile merges secrets into an existing .env file.
// Existing keys are updated if overwrite is true, otherwise skipped.
// New keys are always appended.
func MergeEnvFile(path string, secrets map[string]string, overwrite bool) (MergeResult, error) {
	if path == "" {
		return MergeResult{}, fmt.Errorf("output path must not be empty")
	}

	existing := make(map[string]string)
	var orderedKeys []string

	if f, err := os.Open(path); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			existing[key] = parts[1]
			orderedKeys = append(orderedKeys, key)
		}
	}

	result := MergeResult{}

	// Determine what to update vs skip
	for k, v := range secrets {
		if _, exists := existing[k]; exists {
			if overwrite {
				existing[k] = quoteValue(v)
				result.Updated = append(result.Updated, k)
			} else {
				result.Skipped = append(result.Skipped, k)
			}
		} else {
			orderedKeys = append(orderedKeys, k)
			existing[k] = quoteValue(v)
			result.Added = append(result.Added, k)
		}
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return MergeResult{}, fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, k := range orderedKeys {
		fmt.Fprintf(w, "%s=%s\n", k, existing[k])
	}

	if err := w.Flush(); err != nil {
		return MergeResult{}, fmt.Errorf("failed to flush file: %w", err)
	}

	return result, nil
}
