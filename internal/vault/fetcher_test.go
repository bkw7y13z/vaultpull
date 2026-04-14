package vault

import (
	"context"
	"testing"

	vaultapi "github.com/hashicorp/vault/api"
)

// logicalStub is a minimal stand-in for the Vault logical client used in tests.
type logicalStub struct {
	data   map[string]interface{}
	reaErr error
}

func (l *logicalStub) ReadWithContext(_ context.Context, _ string) (*vaultapi.Secret, error) {
	if l.readErr != nil {
		return nil, l.readErr
	}
	if l.data == nil {
		return nil, nil
	}
	return &vaultapi.Secret{Data: l.data}, nil
}

func newTestClient(stub *logicalStub) *Client {
	return &Client{logical: stub}
}

func TestFetchSecrets_EmptyPath(t *testing.T) {
	c := newTestClient(&logicalStub{})
	_, err := c.FetchSecrets(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestFetchSecrets_VaultError(t *testing.T) {
	stub := &logicalStub{readErr: fmt.Errorf("connection refused")}
	c := newTestClient(stub)
	_, err := c.FetchSecrets(context.Background(), "myapp/prod")
	if err == nil {
		t.Fatal("expected error from vault read, got nil")
	}
}

func TestFetchSecrets_NilSecret(t *testing.T) {
	c := newTestClient(&logicalStub{data: nil})
	_, err := c.FetchSecrets(context.Background(), "myapp/prod")
	if err == nil {
		t.Fatal("expected error for nil secret, got nil")
	}
}

func TestFetchSecrets_Success(t *testing.T) {
	stub := &logicalStub{
		data: map[string]interface{}{
			"data": map[string]interface{}{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
				"IGNORED": 42, // non-string, should be skipped
			},
		},
	}
	c := newTestClient(stub)
	got, err := c.FetchSecrets(context.Background(), "/myapp/prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: want %q, got %q", "localhost", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("DB_PORT: want %q, got %q", "5432", got["DB_PORT"])
	}
	if _, exists := got["IGNORED"]; exists {
		t.Error("non-string value IGNORED should have been skipped")
	}
}
