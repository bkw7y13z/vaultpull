package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretFetcher defines the interface for fetching secrets from Vault.
type SecretFetcher interface {
	FetchSecrets(ctx context.Context, path string) (map[string]string, error)
}

// FetchSecrets reads a KV v2 secret at the given path and returns a flat
// map of string key→value pairs. Only string-typed values are included;
// non-string values are silently skipped.
func (c *Client) FetchSecrets(ctx context.Context, path string) (map[string]string, error) {
	if path == "" {
		return nil, fmt.Errorf("vault: secret path must not be empty")
	}

	// Normalise path: strip leading slash so the Vault SDK receives a
	// clean relative path.
	path = strings.TrimPrefix(path, "/")

	secret, err := c.logical.ReadWithContext(ctx, "secret/data/"+path)
	if err != nil {
		return nil, fmt.Errorf("vault: read %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("vault: no data found at path %q", path)
	}

	// KV v2 wraps the actual values under the "data" key.
	raw, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("vault: unexpected response shape at path %q", path)
	}

	nestedMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("vault: data field is not a map at path %q", path)
	}

	result := make(map[string]string, len(nestedMap))
	for k, v := range nestedMap {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}

	return result, nil
}
