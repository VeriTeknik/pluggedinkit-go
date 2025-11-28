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

// List retrieves all clipboard entries
func (s *ClipboardService) List(ctx context.Context) ([]ClipboardEntry, error) {
	var response ClipboardListResponse
	err := s.client.get(ctx, "/api/clipboard", &response)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, errors.New("failed to list clipboard entries")
	}

	return response.Entries, nil
}

// Get retrieves a clipboard entry by name or index
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
		return nil, nil // Entry not found
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

	// Set defaults
	body := map[string]interface{}{
		"name":  req.Name,
		"value": req.Value,
	}

	if req.ContentType != "" {
		body["contentType"] = req.ContentType
	} else {
		body["contentType"] = "text/plain"
	}

	if req.Encoding != "" {
		body["encoding"] = req.Encoding
	} else {
		body["encoding"] = "utf-8"
	}

	if req.Visibility != "" {
		body["visibility"] = req.Visibility
	} else {
		body["visibility"] = "private"
	}

	if req.CreatedByTool != "" {
		body["createdByTool"] = req.CreatedByTool
	}

	if req.CreatedByModel != "" {
		body["createdByModel"] = req.CreatedByModel
	}

	if req.TTLSeconds > 0 {
		body["ttlSeconds"] = req.TTLSeconds
	}

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

	// Set defaults
	body := map[string]interface{}{
		"value": req.Value,
	}

	if req.ContentType != "" {
		body["contentType"] = req.ContentType
	} else {
		body["contentType"] = "text/plain"
	}

	if req.Encoding != "" {
		body["encoding"] = req.Encoding
	} else {
		body["encoding"] = "utf-8"
	}

	if req.Visibility != "" {
		body["visibility"] = req.Visibility
	} else {
		body["visibility"] = "private"
	}

	if req.CreatedByTool != "" {
		body["createdByTool"] = req.CreatedByTool
	}

	if req.CreatedByModel != "" {
		body["createdByModel"] = req.CreatedByModel
	}

	if req.TTLSeconds > 0 {
		body["ttlSeconds"] = req.TTLSeconds
	}

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

// Pop pops the last indexed entry from clipboard
func (s *ClipboardService) Pop(ctx context.Context) (*ClipboardEntry, error) {
	var response ClipboardResponse
	err := s.client.post(ctx, "/api/clipboard/pop", nil, &response)
	if err != nil {
		return nil, err
	}

	if !response.Success || response.Entry == nil {
		return nil, nil // No entry to pop
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

	return response.Success && response.Deleted, nil
}

// DeleteByName deletes a clipboard entry by name (convenience method)
func (s *ClipboardService) DeleteByName(ctx context.Context, name string) (bool, error) {
	return s.Delete(ctx, &ClipboardDeleteRequest{Name: &name})
}

// DeleteByIndex deletes a clipboard entry by index (convenience method)
func (s *ClipboardService) DeleteByIndex(ctx context.Context, idx int) (bool, error) {
	return s.Delete(ctx, &ClipboardDeleteRequest{Idx: &idx})
}

// ClearAll clears all clipboard entries. Returns count of deleted entries.
func (s *ClipboardService) ClearAll(ctx context.Context) (int, error) {
	entries, err := s.List(ctx)
	if err != nil {
		return 0, err
	}

	deleted := 0
	for _, entry := range entries {
		var success bool
		if entry.Name != nil {
			success, _ = s.DeleteByName(ctx, *entry.Name)
		} else if entry.Idx != nil {
			success, _ = s.DeleteByIndex(ctx, *entry.Idx)
		}
		if success {
			deleted++
		}
	}

	return deleted, nil
}
