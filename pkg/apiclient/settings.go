package apiclient

import (
	"fmt"
	"net/http"
)

// Settings represents the Sitecore Edge admin settings
type Settings struct {
	ContentCacheTtl       string `json:"contentCacheTtl"`
	ContentCacheAutoClear bool   `json:"contentCacheAutoClear"`
	MediaCacheTtl         string `json:"mediaCacheTtl"`
	MediaCacheAutoClear   bool   `json:"mediaCacheAutoClear"`
	TenantCacheAutoClear  bool   `json:"tenantCacheAutoClear"`
}

// NullableSettings - fields are pointers so we can detect which ones are set
type NullableSettings struct {
	ContentCacheTtl       *string `json:"contentCacheTtl,omitempty"`
	ContentCacheAutoClear *bool   `json:"contentCacheAutoClear,omitempty"`
	MediaCacheTtl         *string `json:"mediaCacheTtl,omitempty"`
	MediaCacheAutoClear   *bool   `json:"mediaCacheAutoClear,omitempty"`
	TenantCacheAutoClear  *bool   `json:"tenantCacheAutoClear,omitempty"`
}

// PatchOperation represents a JSON Patch operation (RFC 6902)
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// ToPatchOperations converts only the set fields to PatchOperations
func (s *NullableSettings) ToPatchOperations() []PatchOperation {
	var patches []PatchOperation

	if s.ContentCacheTtl != nil {
		patches = append(patches, PatchOperation{
			Op:    "replace",
			Path:  "/contentcachettl",
			Value: *s.ContentCacheTtl,
		})
	}
	if s.ContentCacheAutoClear != nil {
		patches = append(patches, PatchOperation{
			Op:    "replace",
			Path:  "/contentcacheautoclear",
			Value: *s.ContentCacheAutoClear,
		})
	}
	if s.MediaCacheTtl != nil {
		patches = append(patches, PatchOperation{
			Op:    "replace",
			Path:  "/mediaCacheTtl",
			Value: *s.MediaCacheTtl,
		})
	}
	if s.MediaCacheAutoClear != nil {
		patches = append(patches, PatchOperation{
			Op:    "replace",
			Path:  "/mediaCacheAutoClear",
			Value: *s.MediaCacheAutoClear,
		})
	}
	if s.TenantCacheAutoClear != nil {
		patches = append(patches, PatchOperation{
			Op:    "replace",
			Path:  "/tenantCacheAutoClear",
			Value: *s.TenantCacheAutoClear,
		})
	}

	return patches
}

func (c *Client) GetSettings() (*Settings, error) {
	var settings Settings
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodGet,
		Path:   "/api/admin/v1/settings",
	}, &settings)

	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %v", err)
	}

	return &settings, nil
}

func (c *Client) UpdateSettings(settings *Settings) (*Settings, error) {
	var result Settings
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodPut,
		Path:   "/api/admin/v1/settings",
		Body:   settings,
	}, &result)

	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %v", err)
	}

	return &result, nil
}

func (c *Client) PatchSettings(patches []PatchOperation) (*Settings, error) {
	var result Settings
	err := c.doRequestAndParse(RequestOptions{
		Method: http.MethodPatch,
		Path:   "/api/admin/v1/settings",
		Body:   patches,
	}, &result)

	if err != nil {
		return nil, fmt.Errorf("failed to patch settings: %v", err)
	}

	return &result, nil
}
