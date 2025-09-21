package pluggedinkit

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// DocumentsService handles document-related operations
type DocumentsService struct {
	client *Client
}

// List retrieves a list of documents with optional filters
func (s *DocumentsService) List(ctx context.Context, filters *DocumentFilters) (*ListDocumentsResponse, error) {
	path := "/api/documents"

	if filters != nil {
		params := url.Values{}
		if filters.Source != "" {
			params.Set("source", string(filters.Source))
		}
		if len(filters.Tags) > 0 {
			for _, tag := range filters.Tags {
				params.Add("tags", tag)
			}
		}
		if filters.Category != "" {
			params.Set("category", string(filters.Category))
		}
		if filters.DateFrom != nil {
			params.Set("dateFrom", filters.DateFrom.Format("2006-01-02T15:04:05Z"))
		}
		if filters.DateTo != nil {
			params.Set("dateTo", filters.DateTo.Format("2006-01-02T15:04:05Z"))
		}
		if filters.ModelProvider != "" {
			params.Set("modelProvider", filters.ModelProvider)
		}
		if filters.ModelName != "" {
			params.Set("modelName", filters.ModelName)
		}
		if filters.SearchQuery != "" {
			params.Set("searchQuery", filters.SearchQuery)
		}
		if filters.Sort != "" {
			params.Set("sort", string(filters.Sort))
		}
		if filters.Limit > 0 {
			params.Set("limit", strconv.Itoa(filters.Limit))
		}
		if filters.Offset > 0 {
			params.Set("offset", strconv.Itoa(filters.Offset))
		}

		if len(params) > 0 {
			path = fmt.Sprintf("%s?%s", path, params.Encode())
		}
	}

	var response ListDocumentsResponse
	err := s.client.get(ctx, path, &response)
	return &response, err
}

// Get retrieves a specific document by ID
func (s *DocumentsService) Get(ctx context.Context, documentID string, includeContent, includeVersions bool) (*Document, error) {
	params := url.Values{}
	if includeContent {
		params.Set("includeContent", "true")
	}
	if includeVersions {
		params.Set("includeVersions", "true")
	}

	path := fmt.Sprintf("/api/documents/%s", documentID)
	if len(params) > 0 {
		path = fmt.Sprintf("%s?%s", path, params.Encode())
	}

	var doc Document
	err := s.client.get(ctx, path, &doc)
	return &doc, err
}

// Search performs a semantic search on documents
func (s *DocumentsService) Search(ctx context.Context, query string, filters map[string]interface{}, limit, offset int) (*SearchResponse, error) {
	path := "/api/documents/search"

	params := url.Values{}
	params.Set("query", query)
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	// Add filter parameters
	for key, value := range filters {
		params.Set(key, fmt.Sprintf("%v", value))
	}

	path = fmt.Sprintf("%s?%s", path, params.Encode())

	var response SearchResponse
	err := s.client.get(ctx, path, &response)
	return &response, err
}

// Create creates a new AI-generated document
func (s *DocumentsService) Create(ctx context.Context, req *CreateDocumentRequest) (*Document, error) {
	var doc Document
	err := s.client.post(ctx, "/api/documents", req, &doc)
	return &doc, err
}

// Update updates an existing document
func (s *DocumentsService) Update(ctx context.Context, documentID string, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	path := fmt.Sprintf("/api/documents/%s", documentID)
	var response UpdateDocumentResponse
	err := s.client.put(ctx, path, req, &response)
	return &response, err
}

// Delete deletes a document
func (s *DocumentsService) Delete(ctx context.Context, documentID string) error {
	path := fmt.Sprintf("/api/documents/%s", documentID)
	return s.client.delete(ctx, path, nil)
}