package vault

import (
	"strings"
)

// FilterOptions controls how secrets are filtered before writing.
type FilterOptions struct {
	// Prefix keeps only keys that start with the given prefix.
	Prefix string
	// ExcludeKeys is a set of exact key names to exclude.
	ExcludeKeys []string
}

// FilterSecrets applies FilterOptions to a map of secrets and returns
// a new filtered map. Keys are normalized to uppercase.
func FilterSecrets(secrets map[string]string, opts FilterOptions) map[string]string {
	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[strings.ToUpper(k)] = true
	}

	result := make(map[string]string)
	for k, v := range secrets {
		normKey := strings.ToUpper(k)

		if excluded[normKey] {
			continue
		}

		if opts.Prefix != "" && !strings.HasPrefix(normKey, strings.ToUpper(opts.Prefix)) {
			continue
		}

		result[normKey] = v
	}
	return result
}
