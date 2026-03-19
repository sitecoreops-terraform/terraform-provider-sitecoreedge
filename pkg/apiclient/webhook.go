package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Webhook struct {
	ID            string            `json:"id"`
	TenantID      string            `json:"tenantId"`
	Label         string            `json:"label"`
	URI           string            `json:"uri"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	CreatedBy     string            `json:"createdBy"`
	Created       string            `json:"created"`
	BodyInclude   map[string]string `json:"bodyInclude,omitempty"`
	ExecutionMode string            `json:"executionMode"`
	LastRuns      []interface{}     `json:"lastRuns"`
	Disabled      *bool             `json:"disabled,omitempty"`
}

type WebhookInput struct {
	Label         string            `json:"label"`
	URI           string            `json:"uri"`
	Method        string            `json:"method"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"-"`
	CreatedBy     string            `json:"createdBy"`
	ExecutionMode string            `json:"executionMode"`
	BodyInclude   *string           `json:"bodyInclude,omitempty"`
}

func (w *WebhookInput) MarshalJSON() ([]byte, error) {
	type WebhookInputAlias WebhookInput
	aux := &struct {
		Body        string      `json:"body,omitempty"`
		BodyInclude interface{} `json:"bodyInclude,omitempty"`
		*WebhookInputAlias
	}{
		WebhookInputAlias: (*WebhookInputAlias)(w),
	}

	if w.ExecutionMode != "OnUpdate" && w.Body != "" {
		aux.Body = w.Body
	}

	if w.BodyInclude != nil && *w.BodyInclude != "" {
		var bodyIncludeObj interface{}
		if err := json.Unmarshal([]byte(*w.BodyInclude), &bodyIncludeObj); err == nil {
			aux.BodyInclude = bodyIncludeObj
		} else {
			aux.BodyInclude = *w.BodyInclude
		}
	}

	return json.Marshal(aux)
}

func (c *Client) GetWebhooks() ([]Webhook, error) {
	var webhooks []Webhook
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodGet,
		Path:   "/api/admin/v1/webhooks",
	}, &webhooks)

	if err != nil {
		return nil, fmt.Errorf("failed to get webhooks: %v", err)
	}

	return webhooks, nil
}

func (c *Client) GetWebhook(id string) (*Webhook, error) {
	if id == "" {
		return nil, fmt.Errorf("cannot get webhook, id is required")
	}

	var webhook Webhook
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/admin/v1/webhooks/%s", id),
	}, &webhook)

	if err != nil {
		return nil, fmt.Errorf("failed to get webhook: %v", err)
	}

	return &webhook, nil
}

func (c *Client) CreateWebhook(input *WebhookInput) (*Webhook, error) {
	var webhook Webhook
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodPost,
		Path:   "/api/admin/v1/webhooks",
		Body:   input,
	}, &webhook)

	if err != nil {
		return nil, fmt.Errorf("failed to create webhook: %v", err)
	}

	return &webhook, nil
}

func (c *Client) UpdateWebhook(id string, input *WebhookInput) (*Webhook, error) {

	if id == "" {
		return nil, fmt.Errorf("cannot update webhook, id is necessary")
	}

	_, err := c.doRequest(RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/api/admin/v1/webhooks/%s", id),
		Body:   input,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update webhook: %v", err)
	}

	// Since PUT returns 204 No Content, we need to fetch the updated webhook
	updatedWebhook, err := c.GetWebhook(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated webhook after update: %v", err)
	}

	return updatedWebhook, nil
}

func (c *Client) DeleteWebhook(id string) error {
	_, err := c.doRequest(RequestOptions{
		Method: http.MethodDelete,
		Path:   fmt.Sprintf("/api/admin/v1/webhooks/%s", id),
	})

	if err != nil {
		return fmt.Errorf("failed to delete webhook: %v", err)
	}

	return nil
}
