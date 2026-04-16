package snapshot

import (
	"errors"
	"strings"
)

// RedactOptions controls which keys are redacted in a snapshot entry.
type RedactOptions struct {
	KeyPatterns []string // substrings to match against key names (case-insensitive)
	Replacement string   // value to substitute; defaults to "***REDACTED***"
}

// RedactResult holds the outcome of a redaction pass.
type RedactResult struct {
	RedactedKeys []string
}

// Redact replaces secret values for matching keys in the latest snapshot entry
// and saves the updated snapshot back to path.
func Redact(path string, opts RedactOptions) (*RedactResult, error) {
	if path == "" {
		return nil, errors.New("snapshot path is required")
	}
	if len(opts.KeyPatterns) == 0 {
		return nil, errors.New("at least one key pattern is required")
	}
	replacement := opts.Replacement
	if replacement == "" {
		replacement = "***REDACTED***"
	}

	store, err := Load(path)
	if err != nil {
		return nil, err
	}
	if len(store.Entries) == 0 {
		return nil, errors.New("snapshot is empty")
	}

	latest := &store.Entries[len(store.Entries)-1]
	result := &RedactResult{}

	for k, v := range latest.Secrets {
		for _, pat := range opts.KeyPatterns {
			if strings.Contains(strings.ToLower(k), strings.ToLower(pat)) {
				_ = v
				latest.Secrets[k] = replacement
				result.RedactedKeys = append(result.RedactedKeys, k)
				break
			}
		}
	}

	if err := Save(path, *store); err != nil {
		return nil, err
	}
	return result, nil
}
