package vault

import (
	"errors"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the HashiCorp Vault API client.
type Client struct {
	vc        *vaultapi.Client
	namespace string
}

// Config holds configuration for creating a Vault client.
type Config struct {
	Address   string
	Token     string
	Namespace string
}

// NewClient creates and configures a new Vault client.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Address == "" {
		cfg.Address = os.Getenv("VAULT_ADDR")
	}
	if cfg.Token == "" {
		cfg.Token = os.Getenv("VAULT_TOKEN")
	}
	if cfg.Address == "" {
		return nil, errors.New("vault address is required (set VAULT_ADDR or --vault-addr flag)")
	}
	if cfg.Token == "" {
		return nil, errors.New("vault token is required (set VAULT_TOKEN or --vault-token flag)")
	}

	vcfg := vaultapi.DefaultConfig()
	vcfg.Address = cfg.Address

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}
	vc.SetToken(cfg.Token)

	if cfg.Namespace != "" {
		vc.SetNamespace(cfg.Namespace)
	}

	return &Client{vc: vc, namespace: cfg.Namespace}, nil
}

// ReadSecrets reads key-value secrets from the given KV v2 path.
func (c *Client) ReadSecrets(path string) (map[string]string, error) {
	secret, err := c.vc.KVv2("secret").Get(nil, path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at path %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path %q", path)
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		result[k] = str
	}
	return result, nil
}

// Namespace returns the configured namespace.
func (c *Client) Namespace() string {
	return c.namespace
}
