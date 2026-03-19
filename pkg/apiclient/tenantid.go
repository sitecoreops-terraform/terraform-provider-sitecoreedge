package apiclient

import (
	"fmt"
)

func (c *Client) GetTenantID() (string, error) {
	webhooks, err := c.GetWebhooks()
	if err != nil {
		return "", fmt.Errorf("failed to get webhooks: %v", err)
	}

	if len(webhooks) > 0 {
		return webhooks[0].TenantID, nil
	}

	tempWebhook, err := c.CreateWebhook(&WebhookInput{
		Label:         "temp-tenant-check",
		URI:           "https://temp.example.com",
		Method:        "POST",
		Headers:       map[string]string{},
		Body:          "{}",
		CreatedBy:     "terraform",
		ExecutionMode: "OnEnd",
	})
	if err != nil {
		return "", fmt.Errorf("failed to create temp webhook: %v", err)
	}

	tenantID := tempWebhook.TenantID

	err = c.DeleteWebhook(tempWebhook.ID)
	if err != nil {
		return "", fmt.Errorf("failed to delete temp webhook: %v", err)
	}

	return tenantID, nil
}
