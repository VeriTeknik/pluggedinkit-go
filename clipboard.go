package pluggedinkit

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// ClipboardService handles clipboard-related operations
type ClipboardService struct {
	client *Client
}

// clipboardBodyOptions contains optional fields for building clipboard request bodies
type clipboardBodyOptions struct {
	Name           string
	ContentType    string
	Encoding       ClipboardEncoding
	Visibility     ClipboardVisibility
	CreatedByTool  string
	CreatedByModel string
	TTLSeconds     int
}

// buildClipboardBody builds a request body with proper defaults for clipboard operations
func buildClipboardBody(value string, opts clipboardBodyOptions) map[string]interface{} {
	body := map[string]interface{}{
		"value": value,
	}

	if opts.Name != "" {
		body["name"] = opts.Name
	}

	// Apply defaults for content type
	if opts.ContentType != "" {
		body["contentType"] = opts.ContentType
	} else {
		body["contentType"] = "text/plain"
	}

	// Apply defaults for encoding
	if opts.Encoding != "" {
		body["encoding"] = string(opts.Encoding)
	} else {
		body["encoding"] = string(EncodingUTF8)
	}

	// Apply defaults for visibility
	if opts.Visibility != "" {
		body["visibility"] = string(opts.Visibility)
	} else {
		body["visibility"] = string(ClipboardVisibilityPrivate)
	}

	// Optional fields
	if opts.CreatedByTool != "" {
		body["createdByTool"] = opts.CreatedByTool
	}

	if opts.CreatedByModel != "" {
		body["createdByModel"] = opts.CreatedByModel
	}

	if opts.TTLSeconds > 0 {
		body["ttlSeconds"] = opts.TTLSeconds
	}

	// Hardcode source: SDK always uses 'sdk' source
	body["source"] = string(ClipboardSourceSDK)

	return body
}

// List retrieves clipboard entries with optional pagination.
// limit: maximum number of entries to return (1-100, default: 50)
// offset: number of entries to skip (default: 0)
func (s *ClipboardService) List(ctx context.Context, limit, offset int) (*ClipboardListResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := "/api/clipboard"
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	var response ClipboardListResponse
	err := s.client.get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New("failed to list clipboard entries")
	}

	return &response, nil
}

// Get retrieves a clipboard entry by name or index.
// Returns (nil, nil) if the entry is not found.
// Returns (nil, error) if there was an actual error.
func (s *ClipboardService) Get(ctx context.Context, filters *ClipboardGetFilters) (*ClipboardEntry, error) {
	if filters == nil || (filters.Name == nil && filters.Idx == nil) {
		return nil, errors.New("either 'name' or 'idx' must be provided")
	}

	params := url.Values{}
	if filters.Name != nil {
		params.Set("name", *filters.Name)
	}
	if filters.Idx != nil {
		params.Set("idx", strconv.Itoa(*filters.Idx))
	}

	path := fmt.Sprintf("/api/clipboard?%s", params.Encode())

	var response ClipboardResponse
	err := s.client.get(ctx, path, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		// Distinguish between not-found and actual error
		if response.Error != "" {
			return nil, errors.New(response.Error)
		}
		return nil, nil // Entry not found (no error message)
	}

	return response.Entry, nil
}

// GetByName retrieves a clipboard entry by name (convenience method)
func (s *ClipboardService) GetByName(ctx context.Context, name string) (*ClipboardEntry, error) {
	return s.Get(ctx, &ClipboardGetFilters{Name: &name})
}

// GetByIndex retrieves a clipboard entry by index (convenience method)
func (s *ClipboardService) GetByIndex(ctx context.Context, idx int) (*ClipboardEntry, error) {
	return s.Get(ctx, &ClipboardGetFilters{Idx: &idx})
}

// Set sets a named clipboard entry (upsert)
func (s *ClipboardService) Set(ctx context.Context, req *ClipboardSetRequest) (*ClipboardEntry, error) {
	if req == nil || req.Name == "" {
		return nil, errors.New("name is required")
	}

	body := buildClipboardBody(req.Value, clipboardBodyOptions{
		Name:           req.Name,
		ContentType:    req.ContentType,
		Encoding:       req.Encoding,
		Visibility:     req.Visibility,
		CreatedByTool:  req.CreatedByTool,
		CreatedByModel: req.CreatedByModel,
		TTLSeconds:     req.TTLSeconds,
	})

	var response ClipboardResponse
	err := s.client.post(ctx, "/api/clipboard", body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		if response.Error != "" {
			return nil, errors.New(response.Error)
		}
		return nil, errors.New("failed to set clipboard entry")
	}

	return response.Entry, nil
}

// Push pushes a value to the indexed clipboard (auto-increment index)
func (s *ClipboardService) Push(ctx context.Context, req *ClipboardPushRequest) (*ClipboardEntry, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}

	body := buildClipboardBody(req.Value, clipboardBodyOptions{
		ContentType:    req.ContentType,
		Encoding:       req.Encoding,
		Visibility:     req.Visibility,
		CreatedByTool:  req.CreatedByTool,
		CreatedByModel: req.CreatedByModel,
		TTLSeconds:     req.TTLSeconds,
	})

	var response ClipboardResponse
	err := s.client.post(ctx, "/api/clipboard/push", body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		if response.Error != "" {
			return nil, errors.New(response.Error)
		}
		return nil, errors.New("failed to push to clipboard")
	}

	return response.Entry, nil
}

// Pop pops the last indexed entry from clipboard.
// Returns (nil, nil) if the clipboard is empty.
// Returns (nil, error) if there was an actual error.
func (s *ClipboardService) Pop(ctx context.Context) (*ClipboardEntry, error) {
	var response ClipboardResponse
	err := s.client.post(ctx, "/api/clipboard/pop", nil, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		// Distinguish between empty clipboard and actual error
		if response.Error != "" {
			return nil, errors.New(response.Error)
		}
		return nil, nil // Clipboard is empty (no error message)
	}

	if response.Entry == nil {
		return nil, nil // No entry returned
	}

	return response.Entry, nil
}

// Delete deletes a clipboard entry by name or index
func (s *ClipboardService) Delete(ctx context.Context, req *ClipboardDeleteRequest) (bool, error) {
	if req == nil || (req.Name == nil && req.Idx == nil) {
		return false, errors.New("either 'name' or 'idx' must be provided")
	}

	body := make(map[string]interface{})
	if req.Name != nil {
		body["name"] = *req.Name
	}
	if req.Idx != nil {
		body["idx"] = *req.Idx
	}

	var response ClipboardDeleteResponse
	err := s.client.request(ctx, "DELETE", "/api/clipboard", body, &response)
	if err != nil {
		return false, err
	}

	if !response.Success {
		if response.Error != "" {
			return false, errors.New(response.Error)
		}
		return false, errors.New("failed to delete clipboard entry")
	}

	return response.Deleted, nil
}

// DeleteByName deletes a clipboard entry by name (convenience method)
func (s *ClipboardService) DeleteByName(ctx context.Context, name string) (bool, error) {
	return s.Delete(ctx, &ClipboardDeleteRequest{Name: &name})
}

// DeleteByIndex deletes a clipboard entry by index (convenience method)
func (s *ClipboardService) DeleteByIndex(ctx context.Context, idx int) (bool, error) {
	return s.Delete(ctx, &ClipboardDeleteRequest{Idx: &idx})
}

// ClearAll clears all clipboard entries using bulk delete API.
// Returns a ClearAllResult with the count of deleted entries and success status.
func (s *ClipboardService) ClearAll(ctx context.Context) (*ClearAllResult, error) {
	body := map[string]interface{}{
		"clearAll": true,
	}

	var response struct {
		Success      bool   `json:"success"`
		Deleted      bool   `json:"deleted,omitempty"`
		DeletedCount int    `json:"deletedCount,omitempty"`
		Error        string `json:"error,omitempty"`
	}

	err := s.client.request(ctx, "DELETE", "/api/clipboard", body, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		if response.Error != "" {
			return nil, errors.New(response.Error)
		}
		return nil, errors.New("failed to clear all clipboard entries")
	}

	return &ClearAllResult{
		Deleted: response.DeletedCount,
		Success: true,
	}, nil
}
