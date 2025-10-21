package pluggedinkit

import (
	"context"
	"fmt"
	"net/url"
)

// RAGService handles RAG-related operations
type RAGService struct {
	client *Client
}

// Query performs a RAG query and returns the raw response.
func (s *RAGService) Query(ctx context.Context, query string) (*RAGResponse, error) {
	payload := map[string]interface{}{
		"query":           query,
		"includeMetadata": true,
	}

	var response RAGResponse
	if err := s.client.post(ctx, "/api/rag/query", payload, &response); err != nil {
		return nil, err
	}

	if !response.Success {
		if response.Error != "" {
			return nil, fmt.Errorf("rag query failed: %s", response.Error)
		}
		return nil, fmt.Errorf("rag query failed")
	}

	return &response, nil
}

// AskQuestion performs a simple RAG query and returns the answer text.
func (s *RAGService) AskQuestion(ctx context.Context, query string) (string, error) {
	response, err := s.Query(ctx, query)
	if err != nil {
		return "", err
	}

	if response.Answer == "" {
		return "", fmt.Errorf("no answer returned from knowledge base")
	}

	return response.Answer, nil
}

// QueryWithSources performs a RAG query and returns the response with metadata.
func (s *RAGService) QueryWithSources(ctx context.Context, query, _ string) (*RAGResponse, error) {
	return s.Query(ctx, query)
}

// FindRelevantDocuments finds documents relevant to a query.
func (s *RAGService) FindRelevantDocuments(ctx context.Context, query, _ string, limit int) ([]RAGDocumentReference, error) {
	response, err := s.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	references := make([]RAGDocumentReference, 0, len(response.DocumentIDs))
	for index, documentID := range response.DocumentIDs {
		ref := RAGDocumentReference{
			DocumentID: documentID,
		}
		if index < len(response.Sources) {
			ref.Source = response.Sources[index]
		}
		references = append(references, ref)
	}

	if limit > 0 && len(references) > limit {
		return references[:limit], nil
	}

	return references, nil
}

// CheckAvailability checks if the RAG service is reachable.
func (s *RAGService) CheckAvailability(ctx context.Context) (map[string]interface{}, error) {
	_, err := s.Query(ctx, "__pluggedin_health_check__")
	if err != nil {
		return map[string]interface{}{
			"available": false,
			"message":   err.Error(),
		}, nil
	}

	return map[string]interface{}{
		"available": true,
	}, nil
}

// GetStorageStats gets storage statistics for RAG.
func (s *RAGService) GetStorageStats(ctx context.Context, userID string) (*RAGStorageStats, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is required to fetch storage statistics")
	}

	params := url.Values{}
	params.Set("user_id", userID)

	path := fmt.Sprintf("/api/rag/storage-stats?%s", params.Encode())

	var stats RAGStorageStats
	if err := s.client.get(ctx, path, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
