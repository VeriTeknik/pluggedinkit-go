package pluggedinkit

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// RAGService handles RAG-related operations
type RAGService struct {
	client *Client
}

// AskQuestion performs a simple RAG query
func (s *RAGService) AskQuestion(ctx context.Context, query string) (string, error) {
	request := map[string]string{
		"query": query,
	}

	var response map[string]interface{}
	err := s.client.post(ctx, "/api/library/rag/query", request, &response)
	if err != nil {
		return "", err
	}

	answer, ok := response["answer"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format")
	}

	return answer, nil
}

// QueryWithSources performs a RAG query and returns sources
func (s *RAGService) QueryWithSources(ctx context.Context, query, projectUUID string) (*RAGResponse, error) {
	request := map[string]string{
		"query": query,
	}

	if projectUUID != "" {
		request["projectUuid"] = projectUUID
	}

	var response RAGResponse
	err := s.client.post(ctx, "/api/library/rag/query", request, &response)
	return &response, err
}

// FindRelevantDocuments finds documents relevant to a query
func (s *RAGService) FindRelevantDocuments(ctx context.Context, query, projectUUID string, limit int) ([]RAGDocument, error) {
	params := url.Values{}
	params.Set("query", query)
	if projectUUID != "" {
		params.Set("projectUuid", projectUUID)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	path := fmt.Sprintf("/api/library/rag/relevant?%s", params.Encode())

	var response struct {
		Documents []RAGDocument `json:"documents"`
	}
	err := s.client.get(ctx, path, &response)
	return response.Documents, err
}

// CheckAvailability checks if RAG is available
func (s *RAGService) CheckAvailability(ctx context.Context) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.get(ctx, "/api/library/rag/status", &response)
	return response, err
}

// GetStorageStats gets storage statistics for RAG
func (s *RAGService) GetStorageStats(ctx context.Context) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := s.client.get(ctx, "/api/library/rag/stats", &response)
	return response, err
}