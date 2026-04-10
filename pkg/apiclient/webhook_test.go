package apiclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetWebhooks(t *testing.T) {
	sampleResponse := `[{"id":"03a51bcc-e741-4e10-92d1-f1c9ab41e481","tenantId":"pkaejendomma766-pkaejendomm4e8c-deve3c6-e49a","label":"test","uri":"https://dev-api.pka-ejendomme.dk/api/webhook1","method":"POST","headers":{"x-header":"bar"},"body":"{\"rebuild\":\"true\"}","createdBy":"mig","created":"2026-02-20T23:37:22.2297746+00:00","bodyInclude":{"updated":"true"},"executionMode":"OnUpdate","lastRuns":[],"disabled":null}]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/admin/v1/webhooks" {
			t.Errorf("Expected path /api/admin/v1/webhooks, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(sampleResponse))
		if err != nil {
			t.Fatalf("TestGetWebhooks write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	webhooks, err := client.GetWebhooks()
	if err != nil {
		t.Fatalf("GetWebhooks failed: %v", err)
	}

	if len(webhooks) != 1 {
		t.Fatalf("Expected 1 webhook, got %d", len(webhooks))
	}

	webhook := webhooks[0]
	if webhook.ID != "03a51bcc-e741-4e10-92d1-f1c9ab41e481" {
		t.Errorf("Expected ID '03a51bcc-e741-4e10-92d1-f1c9ab41e481', got '%s'", webhook.ID)
	}
	if webhook.TenantID != "pkaejendomma766-pkaejendomm4e8c-deve3c6-e49a" {
		t.Errorf("Expected TenantID 'pkaejendomma766-pkaejendomm4e8c-deve3c6-e49a', got '%s'", webhook.TenantID)
	}
	if webhook.Label != "test" {
		t.Errorf("Expected Label 'test', got '%s'", webhook.Label)
	}
	if webhook.URI != "https://dev-api.pka-ejendomme.dk/api/webhook1" {
		t.Errorf("Expected URI 'https://dev-api.pka-ejendomme.dk/api/webhook1', got '%s'", webhook.URI)
	}
	if webhook.Method != "POST" {
		t.Errorf("Expected Method 'POST', got '%s'", webhook.Method)
	}
	if webhook.Headers["x-header"] != "bar" {
		t.Errorf("Expected header 'x-header' to be 'bar', got '%s'", webhook.Headers["x-header"])
	}
	if webhook.Body != `{"rebuild":"true"}` {
		t.Errorf("Expected Body '{\"rebuild\":\"true\"}', got '%s'", webhook.Body)
	}
	if webhook.CreatedBy != "mig" {
		t.Errorf("Expected CreatedBy 'mig', got '%s'", webhook.CreatedBy)
	}
	if webhook.Created != "2026-02-20T23:37:22.2297746+00:00" {
		t.Errorf("Expected Created '2026-02-20T23:37:22.2297746+00:00', got '%s'", webhook.Created)
	}
	if webhook.ExecutionMode != "OnUpdate" {
		t.Errorf("Expected ExecutionMode 'OnUpdate', got '%s'", webhook.ExecutionMode)
	}
	if webhook.BodyInclude["updated"] != "true" {
		t.Errorf("Expected BodyInclude[updated] to be 'true', got '%s'", webhook.BodyInclude["updated"])
	}
	if webhook.Disabled != nil {
		t.Errorf("Expected Disabled to be nil, got %v", *webhook.Disabled)
	}
}

func TestGetWebhooksEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("[]"))
		if err != nil {
			t.Fatalf("TestGetWebhooksEmptyResponse write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	webhooks, err := client.GetWebhooks()
	if err != nil {
		t.Fatalf("GetWebhooks failed: %v", err)
	}

	if len(webhooks) != 0 {
		t.Errorf("Expected 0 webhooks, got %d", len(webhooks))
	}
}

func TestGetWebhook(t *testing.T) {
	sampleResponse := `{"id":"test-id","tenantId":"test-tenant","label":"test-label","uri":"https://example.com","method":"POST","headers":{"x-header":"value"},"body":"test-body","createdBy":"test-user","created":"2026-03-18T00:00:00.0000000+00:00","bodyInclude":{"key":"value"},"executionMode":"OnEnd","lastRuns":[],"disabled":false}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/admin/v1/webhooks/test-id") {
			t.Errorf("Expected path to contain /api/admin/v1/webhooks/test-id, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(sampleResponse))
		if err != nil {
			t.Fatalf("TestGetWebhook write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	webhook, err := client.GetWebhook("test-id")
	if err != nil {
		t.Fatalf("GetWebhook failed: %v", err)
	}

	if webhook == nil {
		t.Fatal("Expected webhook to be returned, got nil")
	}

	if webhook.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", webhook.ID)
	}
	if webhook.Label != "test-label" {
		t.Errorf("Expected Label 'test-label', got '%s'", webhook.Label)
	}
	if webhook.URI != "https://example.com" {
		t.Errorf("Expected URI 'https://example.com', got '%s'", webhook.URI)
	}
	if webhook.Method != "POST" {
		t.Errorf("Expected Method 'POST', got '%s'", webhook.Method)
	}
	if webhook.Body != "test-body" {
		t.Errorf("Expected Body 'test-body', got '%s'", webhook.Body)
	}
}

func TestGetWebhookRequiresID(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	_, err := client.GetWebhook("")
	if err == nil {
		t.Fatal("Expected error when ID is empty, got nil")
	}

	if err.Error() != "cannot get webhook, id is required" {
		t.Errorf("Expected error 'cannot get webhook, id is required', got '%s'", err.Error())
	}
}

func TestUpdateWebhookRequiresID(t *testing.T) {
	client := &Client{
		BaseURL:    "https://example.com",
		Token:      "test-token",
		HTTPClient: &http.Client{},
	}

	input := &WebhookInput{
		Label:         "test",
		URI:           "https://example.com",
		Method:        "POST",
		Headers:       map[string]string{},
		CreatedBy:     "test",
		ExecutionMode: "OnEnd",
	}

	_, err := client.UpdateWebhook("", input)
	if err == nil {
		t.Fatal("Expected error when ID is empty, got nil")
	}

	if err.Error() != "cannot update webhook, id is necessary" {
		t.Errorf("Expected error 'cannot update webhook, id is necessary', got '%s'", err.Error())
	}
}

func TestCreateWebhookOnUpdateExcludesBody(t *testing.T) {
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if _, exists := receivedBody["body"]; exists {
			t.Error("Request body should NOT contain 'body' field when ExecutionMode is OnUpdate")
		}

		if receivedBody["executionMode"] != "OnUpdate" {
			t.Errorf("Expected executionMode 'OnUpdate', got '%s'", receivedBody["executionMode"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"id":"test-id","label":"test","uri":"https://example.com","method":"POST","executionMode":"OnUpdate","createdBy":"test"}`))
		if err != nil {
			t.Fatalf("TestCreateWebookOnUpdateExcludesBody write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	input := &WebhookInput{
		Label:         "test",
		URI:           "https://example.com",
		Method:        "POST",
		Headers:       map[string]string{"X-Test": "value"},
		Body:          "should-be-excluded",
		CreatedBy:     "test",
		ExecutionMode: "OnUpdate",
	}

	webhook, err := client.CreateWebhook(input)
	if err != nil {
		t.Fatalf("CreateWebhook failed: %v", err)
	}

	if webhook.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", webhook.ID)
	}
}

func TestCreateWebhookWithBodyInclude(t *testing.T) {
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if _, exists := receivedBody["body"]; exists {
			t.Error("Request body should NOT contain 'body' field when ExecutionMode is OnUpdate")
		}

		bodyInclude, exists := receivedBody["bodyInclude"]
		if !exists {
			t.Fatal("Request body should contain 'bodyInclude' field")
		}

		bodyIncludeMap, ok := bodyInclude.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected bodyInclude to be an object, got %T", bodyInclude)
		}

		if bodyIncludeMap["updated"] != "true" {
			t.Errorf("Expected bodyInclude.updated to be 'true', got '%s'", bodyIncludeMap["updated"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"id":"test-id","label":"test","uri":"https://example.com","method":"POST","bodyInclude":{"updated":"true"},"executionMode":"OnUpdate","createdBy":"test"}`))
		if err != nil {
			t.Fatalf("TestCreateWebhooksWithBodyIncluded write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	bodyIncludeStr := `{"updated":"true"}`
	input := &WebhookInput{
		Label:         "test",
		URI:           "https://example.com",
		Method:        "POST",
		Headers:       map[string]string{"X-Test": "value"},
		Body:          "should-be-excluded",
		CreatedBy:     "test",
		ExecutionMode: "OnUpdate",
		BodyInclude:   &bodyIncludeStr,
	}

	webhook, err := client.CreateWebhook(input)
	if err != nil {
		t.Fatalf("CreateWebhook failed: %v", err)
	}

	if webhook.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", webhook.ID)
	}
}

func TestCreateWebhookOnEndIncludesBody(t *testing.T) {
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		if _, exists := receivedBody["body"]; !exists {
			t.Error("Request body SHOULD contain 'body' field when ExecutionMode is OnEnd")
		}

		if receivedBody["body"] != "expected-body-content" {
			t.Errorf("Expected body 'expected-body-content', got '%s'", receivedBody["body"])
		}

		if receivedBody["executionMode"] != "OnEnd" {
			t.Errorf("Expected executionMode 'OnEnd', got '%s'", receivedBody["executionMode"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"id":"test-id","label":"test","uri":"https://example.com","method":"POST","body":"expected-body-content","executionMode":"OnEnd","createdBy":"test"}`))
		if err != nil {
			t.Fatalf("TestCreateWebhooksOnEndIncludesBody write failed: %v", err)
		}
	}))
	defer server.Close()

	client := &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}

	input := &WebhookInput{
		Label:         "test",
		URI:           "https://example.com",
		Method:        "POST",
		Headers:       map[string]string{"X-Test": "value"},
		Body:          "expected-body-content",
		CreatedBy:     "test",
		ExecutionMode: "OnEnd",
	}

	webhook, err := client.CreateWebhook(input)
	if err != nil {
		t.Fatalf("CreateWebhook failed: %v", err)
	}

	if webhook.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", webhook.ID)
	}
}

func TestWebhookInputSerialization(t *testing.T) {
	tests := []struct {
		name          string
		input         *WebhookInput
		expectBody    bool
		expectBodyVal string
	}{
		{
			name: "OnUpdate mode should exclude body",
			input: &WebhookInput{
				Label:         "test",
				URI:           "https://example.com",
				Method:        "POST",
				Headers:       map[string]string{},
				Body:          "should-not-send",
				CreatedBy:     "test",
				ExecutionMode: "OnUpdate",
			},
			expectBody:    false,
			expectBodyVal: "should-not-send",
		},
		{
			name: "OnEnd mode should include body",
			input: &WebhookInput{
				Label:         "test",
				URI:           "https://example.com",
				Method:        "POST",
				Headers:       map[string]string{},
				Body:          "should-send",
				CreatedBy:     "test",
				ExecutionMode: "OnEnd",
			},
			expectBody:    true,
			expectBodyVal: "should-send",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("Failed to marshal WebhookInput: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &result); err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			_, hasBody := result["body"]
			if tt.expectBody && !hasBody {
				t.Errorf("Expected body field to be present for ExecutionMode %s", tt.input.ExecutionMode)
			}
			if !tt.expectBody && hasBody {
				t.Errorf("Expected body field to be absent for ExecutionMode %s", tt.input.ExecutionMode)
			}
		})
	}
}

func TestWebhookInputJSONOmitsEmptyBodyOnUpdate(t *testing.T) {
	input := &WebhookInput{
		Label:         "test",
		URI:           "https://example.com",
		Method:        "POST",
		Headers:       map[string]string{},
		Body:          "",
		CreatedBy:     "test",
		ExecutionMode: "OnUpdate",
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal WebhookInput: %v", err)
	}

	jsonStr := string(jsonBytes)
	if strings.Contains(jsonStr, `"body"`) {
		t.Error("JSON should not contain 'body' field when empty and ExecutionMode is OnUpdate")
	}
}

func TestUpdateWebhook(t *testing.T) {
	sampleWebhookResponse := `{"id":"test-id","tenantId":"test-tenant","label":"updated-label","uri":"https://updated.example.com","method":"PUT","headers":{"x-custom":"header"},"body":"updated-body","createdBy":"test-user","created":"2026-03-18T00:00:00.0000000+00:00","bodyInclude":{"key":"value"},"executionMode":"OnEnd","lastRuns":[],"disabled":false}`

	var putCallCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api/admin/v1/webhooks/test-id") {
			switch r.Method {
			case http.MethodPut:
				putCallCount++
				var receivedBody map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
					t.Fatalf("Failed to decode request body: %v", err)
				}

				if receivedBody["label"] != "updated-label" {
					t.Errorf("Expected label 'updated-label', got '%s'", receivedBody["label"])
				}

				// PUT returns 204 No Content
				w.WriteHeader(http.StatusNoContent)
				return
			case http.MethodGet:
				// This is the subsequent GET request to fetch the updated webhook
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte(sampleWebhookResponse))
				if err != nil {
					t.Fatalf("TestUpdateWebook write failed: %v", err)
				}
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

	input := &WebhookInput{
		Label:         "updated-label",
		URI:           "https://updated.example.com",
		Method:        "PUT",
		Headers:       map[string]string{"x-custom": "header"},
		Body:          "updated-body",
		CreatedBy:     "test-user",
		ExecutionMode: "OnEnd",
	}

	webhook, err := client.UpdateWebhook("test-id", input)
	if err != nil {
		t.Fatalf("UpdateWebhook failed: %v", err)
	}

	if webhook == nil {
		t.Fatal("Expected webhook to be returned, got nil")
	}

	if putCallCount != 1 {
		t.Errorf("Expected 1 PUT call, got %d", putCallCount)
	}

	if webhook.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", webhook.ID)
	}
	if webhook.Label != "updated-label" {
		t.Errorf("Expected Label 'updated-label', got '%s'", webhook.Label)
	}
	if webhook.URI != "https://updated.example.com" {
		t.Errorf("Expected URI 'https://updated.example.com', got '%s'", webhook.URI)
	}
	if webhook.Method != "PUT" {
		t.Errorf("Expected Method 'PUT', got '%s'", webhook.Method)
	}
	if webhook.Body != "updated-body" {
		t.Errorf("Expected Body 'updated-body', got '%s'", webhook.Body)
	}
}
