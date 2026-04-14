package vault

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_MissingAddress(t *testing.T) {
	t.Setenv("VAULT_ADDR", "")
	t.Setenv("VAULT_TOKEN", "test-token")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when VAULT_ADDR is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when VAULT_TOKEN is missing, got nil")
	}
}

func TestNewClient_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(Config{
		Address:   server.URL,
		Token:     "test-token",
		Namespace: "my-ns",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Namespace() != "my-ns" {
		t.Errorf("expected namespace %q, got %q", "my-ns", client.Namespace())
	}
}

func TestNewClient_EnvFallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	t.Setenv("VAULT_ADDR", server.URL)
	t.Setenv("VAULT_TOKEN", "env-token")

	client, err := NewClient(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
