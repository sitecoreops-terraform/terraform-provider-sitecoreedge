package apiclient

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListTokens(t *testing.T) {
	sampleResponse := `{
		"totalCount": 1,
		"pageSize": 20,
		"currentPage": 1,
		"totalPages": 1,
		"hasNext": false,
		"hasPrevious": false,
		"keys": [
			{
				"tenantId": "test-tenant",
				"hash": "test-hash",
				"isRevoked": false,
				"label": "Test Token",
				"scopes": ["content-#everything#", "audience-delivery"],
				"createdBy": "test-user",
				"created": "2026-04-22T00:00:00.0000000+00:00"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/apikey/v1" {
			t.Errorf("Expected path /api/apikey/v1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(sampleResponse))
		if err != nil {
			t.Fatalf("TestListTokens write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	tokens, err := client.ListTokens(false, "")
	if err != nil {
		t.Fatalf("ListTokens failed: %v", err)
	}

	if len(tokens) != 1 {
		t.Fatalf("Expected 1 token, got %d", len(tokens))
	}

	token := tokens[0]
	if token.Hash != "test-hash" {
		t.Errorf("Expected Hash 'test-hash', got '%s'", token.Hash)
	}
	if token.Label != "Test Token" {
		t.Errorf("Expected Label 'Test Token', got '%s'", token.Label)
	}
	if len(token.Scopes) != 2 {
		t.Errorf("Expected 2 scopes, got %d", len(token.Scopes))
	}
}

func TestGetTokenByHash(t *testing.T) {
	sampleResponse := `{
		"tenantId": "test-tenant",
		"hash": "test-hash",
		"isRevoked": false,
		"label": "Test Token",
		"scopes": ["content-#everything#", "audience-delivery"],
		"createdBy": "test-user",
		"created": "2026-04-22T00:00:00.0000000+00:00"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/apikey/v1/test-hash") {
			t.Errorf("Expected path to contain /api/apikey/v1/test-hash, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(sampleResponse))
		if err != nil {
			t.Fatalf("TestGetTokenByHash write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	token, err := client.GetTokenByHash("test-hash")
	if err != nil {
		t.Fatalf("GetTokenByHash failed: %v", err)
	}

	if token == nil {
		t.Fatal("Expected token to be returned, got nil")
	}

	if token.Hash != "test-hash" {
		t.Errorf("Expected Hash 'test-hash', got '%s'", token.Hash)
	}
	if token.Label != "Test Token" {
		t.Errorf("Expected Label 'Test Token', got '%s'", token.Label)
	}
}

func TestRenameToken(t *testing.T) {
	var putCallCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/apikey/v1/renamebyhash/test-hash") {
			switch r.Method {
			case http.MethodPut:
				putCallCount++
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	err := client.RenameToken("test-hash", "New Token Name")
	if err != nil {
		t.Fatalf("RenameToken failed: %v", err)
	}

	if putCallCount != 1 {
		t.Errorf("Expected 1 PUT call, got %d", putCallCount)
	}
}

func TestRevokeTokenByHash(t *testing.T) {
	var putCallCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/apikey/v1/revokebyhash/test-hash") {
			switch r.Method {
			case http.MethodPut:
				putCallCount++
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	err := client.RevokeTokenByHash("test-hash")
	if err != nil {
		t.Fatalf("RevokeTokenByHash failed: %v", err)
	}

	if putCallCount != 1 {
		t.Errorf("Expected 1 PUT call, got %d", putCallCount)
	}
}

func TestCreateToken(t *testing.T) {
	expectedTokenValue := "test-token-value-12345"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(expectedTokenValue))
		if err != nil {
			t.Fatalf("TestCreateToken write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	input := &TokenInput{
		CreatedBy: "test-user",
		Label:     "Test Token",
		Scopes:    []string{"content-#everything#", "audience-delivery"},
	}

	tokenValue, err := client.CreateToken(input)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	if tokenValue != expectedTokenValue {
		t.Errorf("Expected token value '%s', got '%s'", expectedTokenValue, tokenValue)
	}
}

func TestGetTokenByToken(t *testing.T) {
	sampleResponse := `{
		"tenantId": "test-tenant",
		"hash": "test-hash",
		"isRevoked": false,
		"label": "Test Token",
		"scopes": ["content-#everything#", "audience-delivery"],
		"createdBy": "test-user",
		"created": "2026-04-22T00:00:00.0000000+00:00"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/apikey/v1/token" {
			t.Errorf("Expected path /api/apikey/v1/token, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check for X-GQL-Token header
		tokenHeader := r.Header.Get("X-GQL-Token")
		if tokenHeader != "test-token-value" {
			t.Errorf("Expected X-GQL-Token header 'test-token-value', got '%s'", tokenHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(sampleResponse))
		if err != nil {
			t.Fatalf("TestGetTokenByToken write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	token, err := client.GetTokenByToken("test-token-value")
	if err != nil {
		t.Fatalf("GetTokenByToken failed: %v", err)
	}

	if token == nil {
		t.Fatal("Expected token to be returned, got nil")
	}

	if token.Hash != "test-hash" {
		t.Errorf("Expected Hash 'test-hash', got '%s'", token.Hash)
	}
	if token.Label != "Test Token" {
		t.Errorf("Expected Label 'Test Token', got '%s'", token.Label)
	}
}

func TestRevokeTokenByToken(t *testing.T) {
	var putCallCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/apikey/v1/revokebytoken" {
			switch r.Method {
			case http.MethodPut:
				putCallCount++

				// Check for X-GQL-Token header
				tokenHeader := r.Header.Get("X-GQL-Token")
				if tokenHeader != "test-token-value" {
					t.Errorf("Expected X-GQL-Token header 'test-token-value', got '%s'", tokenHeader)
				}

				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		t.Errorf("Unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	err := client.RevokeTokenByToken("test-token-value")
	if err != nil {
		t.Fatalf("RevokeTokenByToken failed: %v", err)
	}

	if putCallCount != 1 {
		t.Errorf("Expected 1 PUT call, got %d", putCallCount)
	}
}

func TestListTokensRequiresHash(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	_, err := client.GetTokenByHash("")
	if err == nil {
		t.Fatal("Expected error when hash is empty, got nil")
	}

	if err.Error() != "hash is required to get token" {
		t.Errorf("Expected error 'hash is required to get token', got '%s'", err.Error())
	}
}

func TestRenameTokenRequiresHash(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	err := client.RenameToken("", "New Name")
	if err == nil {
		t.Fatal("Expected error when hash is empty, got nil")
	}

	if err.Error() != "hash is required to rename token" {
		t.Errorf("Expected error 'hash is required to rename token', got '%s'", err.Error())
	}
}

func TestRevokeTokenByHashRequiresHash(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	err := client.RevokeTokenByHash("")
	if err == nil {
		t.Fatal("Expected error when hash is empty, got nil")
	}

	if err.Error() != "hash is required to revoke token" {
		t.Errorf("Expected error 'hash is required to revoke token', got '%s'", err.Error())
	}
}

func TestCreateTokenRequiresInput(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	_, err := client.CreateToken(nil)
	if err == nil {
		t.Fatal("Expected error when input is nil, got nil")
	}

	if err.Error() != "input is required to create token" {
		t.Errorf("Expected error 'input is required to create token', got '%s'", err.Error())
	}
}

func TestGetTokenByTokenRequiresToken(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	_, err := client.GetTokenByToken("")
	if err == nil {
		t.Fatal("Expected error when token is empty, got nil")
	}

	if err.Error() != "token value is required" {
		t.Errorf("Expected error 'token value is required', got '%s'", err.Error())
	}
}

func TestRevokeTokenByTokenRequiresToken(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	err := client.RevokeTokenByToken("")
	if err == nil {
		t.Fatal("Expected error when token is empty, got nil")
	}

	if err.Error() != "token value is required to revoke token" {
		t.Errorf("Expected error 'token value is required to revoke token', got '%s'", err.Error())
	}
}
