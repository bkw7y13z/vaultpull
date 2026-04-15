package snapshot

import (
	"fmt"
	"strings"
)

// LintIssue describes a single lint warning or error found in a snapshot entry.
type LintIssue struct {
	Checksum string
	Field    string
	Message  string
}

func (i LintIssue) String() string {
	return fmt.Sprintf("[%s] %s: %s", i.Checksum[:8], i.Field, i.Message)
}

// LintResult holds all issues found during a lint pass.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasIssues() bool {
	return len(r.Issues) > 0
}

// Lint validates the snapshot file at path and returns a LintResult.
// It checks for empty keys, duplicate checksums, missing timestamps, and
// entries with no keys recorded.
func Lint(path string) (*LintResult, error) {
	if path == "" {
		return nil, fmt.Errorf("snapshot path must not be empty")
	}

	entries, err := Load(path)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot: %w", err)
	}

	result := &LintResult{}
	seen := make(map[string]bool)

	for _, e := range entries {
		cs := e.Checksum

		if cs == "" {
			result.Issues = append(result.Issues, LintIssue{
				Checksum: "(empty)",
				Field:    "checksum",
				Message:  "entry has empty checksum",
			})
			continue
		}

		if seen[cs] {
			result.Issues = append(result.Issues, LintIssue{
				Checksum: cs,
				Field:    "checksum",
				Message:  "duplicate checksum detected",
			})
		}
		seen[cs] = true

		if e.CreatedAt.IsZero() {
			result.Issues = append(result.Issues, LintIssue{
				Checksum: cs,
				Field:    "created_at",
				Message:  "missing creation timestamp",
			})
		}

		if len(e.Keys) == 0 {
			result.Issues = append(result.Issues, LintIssue{
				Checksum: cs,
				Field:    "keys",
				Message:  "entry has no keys recorded",
			})
		}

		for _, k := range e.Keys {
			if strings.TrimSpace(k) == "" {
				result.Issues = append(result.Issues, LintIssue{
					Checksum: cs,
					Field:    "keys",
					Message:  "entry contains blank key name",
				})
				break
			}
		}
	}

	return result, nil
}
