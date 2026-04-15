package vault

import (
	"testing"
)

func TestFilterSecrets_NoFilter(t *testing.T) {
	secrets := map[string]string{"DB_HOST": "localhost", "API_KEY": "abc123"}
	result := FilterSecrets(secrets, FilterOptions{})
	if len(result) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(result))
	}
}

func TestFilterSecrets_PrefixFilter(t *testing.T) {
	secrets := map[string]string{
		"db_host":     "localhost",
		"db_password": "secret",
		"api_key":     "abc123",
	}
	result := FilterSecrets(secrets, FilterOptions{Prefix: "DB_"})
	if len(result) != 2 {
		t.Errorf("expected 2 secrets after prefix filter, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := result["API_KEY"]; ok {
		t.Error("expected API_KEY to be filtered out")
	}
}

func TestFilterSecrets_ExcludeKeys(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
	}
	result := FilterSecrets(secrets, FilterOptions{ExcludeKeys: []string{"db_password", "API_KEY"}})
	if len(result) != 1 {
		t.Errorf("expected 1 secret after exclusion, got %d", len(result))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
}

func TestFilterSecrets_NormalizesKeys(t *testing.T) {
	secrets := map[string]string{"db_host": "localhost"}
	result := FilterSecrets(secrets, FilterOptions{})
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected key to be normalized to uppercase DB_HOST")
	}
}

func TestFilterSecrets_PrefixAndExclude(t *testing.T) {
	secrets := map[string]string{
		"APP_NAME":    "myapp",
		"APP_SECRET":  "topsecret",
		"OTHER_VALUE": "ignore",
	}
	result := FilterSecrets(secrets, FilterOptions{
		Prefix:      "APP_",
		ExcludeKeys: []string{"APP_SECRET"},
	})
	if len(result) != 1 {
		t.Errorf("expected 1 secret, got %d", len(result))
	}
	if _, ok := result["APP_NAME"]; !ok {
		t.Error("expected APP_NAME in result")
	}
}

func TestFilterSecrets_EmptyInput(t *testing.T) {
	secrets := map[string]string{}
	result := FilterSecrets(secrets, FilterOptions{Prefix: "APP_"})
	if len(result) != 0 {
		t.Errorf("expected 0 secrets for empty input, got %d", len(result))
	}
}
