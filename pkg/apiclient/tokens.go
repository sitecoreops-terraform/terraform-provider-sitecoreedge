package apiclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Token struct {
	TenantID  string   `json:"tenantId"`
	Hash      string   `json:"hash"`
	IsRevoked bool     `json:"isRevoked"`
	Label     string   `json:"label"`
	Scopes    []string `json:"scopes"`
	CreatedBy string   `json:"createdBy"`
	Created   string   `json:"created"`
}

type TokenInput struct {
	CreatedBy string   `json:"createdBy"`
	Label     string   `json:"label"`
	Scopes    []string `json:"scopes"`
}

type TokenListResponse struct {
	TotalCount  int     `json:"totalCount"`
	PageSize    int     `json:"pageSize"`
	CurrentPage int     `json:"currentPage"`
	TotalPages  int     `json:"totalPages"`
	HasNext     bool    `json:"hasNext"`
	HasPrevious bool    `json:"hasPrevious"`
	Keys        []Token `json:"keys"`
}

type RenameTokenInput struct {
	NewName string `json:"newName"`
}

func (c *Client) ListTokens(filterRevoked bool, label string) ([]Token, error) {
	var response TokenListResponse

	path := "/api/apikey/v1"
	queryParams := []string{}

	if filterRevoked {
		queryParams = append(queryParams, "filterRevoked=true")
	}

	if label != "" {
		queryParams = append(queryParams, fmt.Sprintf("label=%s", label))
	}

	if len(queryParams) > 0 {
		path += "?" + strings.Join(queryParams, "&")
	}

	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodGet,
		Path:   path,
	}, &response)

	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %v", err)
	}

	return response.Keys, nil
}

func (c *Client) GetTokenByHash(hash string) (*Token, error) {
	if hash == "" {
		return nil, fmt.Errorf("hash is required to get token")
	}

	var token Token
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/apikey/v1/%s", hash),
	}, &token)

	if err != nil {
		return nil, fmt.Errorf("failed to get token by hash: %v", err)
	}

	return &token, nil
}

func (c *Client) RenameToken(hash string, newName string) error {
	if hash == "" {
		return fmt.Errorf("hash is required to rename token")
	}

	input := &RenameTokenInput{
		NewName: newName,
	}

	_, err := c.doRequest(RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/api/apikey/v1/renamebyhash/%s", hash),
		Body:   input,
	})

	if err != nil {
		return fmt.Errorf("failed to rename token: %v", err)
	}

	return nil
}

func (c *Client) RevokeTokenByHash(hash string) error {
	if hash == "" {
		return fmt.Errorf("hash is required to revoke token")
	}

	_, err := c.doRequest(RequestOptions{
		Method: http.MethodPut,
		Path:   fmt.Sprintf("/api/apikey/v1/revokebyhash/%s", hash),
	})

	if err != nil {
		return fmt.Errorf("failed to revoke token by hash: %v", err)
	}

	return nil
}

func (c *Client) CreateToken(input *TokenInput) (string, error) {
	if input == nil {
		return "", fmt.Errorf("input is required to create token")
	}

	resp, err := c.doRequest(RequestOptions{
		Method: http.MethodPost,
		Path:   "/api/apikey/v1",
		Body:   input,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create token: %v", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// The response is plain text with the token value
	tokenBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %v", err)
	}

	return string(tokenBytes), nil
}

func (c *Client) GetTokenByToken(tokenValue string) (*Token, error) {
	if tokenValue == "" {
		return nil, fmt.Errorf("token value is required")
	}

	var token Token
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodGet,
		Path:   "/api/apikey/v1/token",
		Headers: map[string]string{
			"X-GQL-Token": tokenValue,
		},
	}, &token)

	if err != nil {
		return nil, fmt.Errorf("failed to get token by token: %v", err)
	}

	return &token, nil
}

func (c *Client) RevokeTokenByToken(tokenValue string) error {
	if tokenValue == "" {
		return fmt.Errorf("token value is required to revoke token")
	}

	_, err := c.doRequest(RequestOptions{
		Method: http.MethodPut,
		Path:   "/api/apikey/v1/revokebytoken",
		Headers: map[string]string{
			"X-GQL-Token": tokenValue,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to revoke token by token: %v", err)
	}

	return nil
}
