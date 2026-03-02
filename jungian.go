package pluggedinkit

import (
	"context"
	"fmt"
	"net/url"
)

// JungianService handles Jungian Intelligence Layer operations
type JungianService struct {
	client *Client
}

// SearchWithContext performs an archetype-enhanced memory search
func (s *JungianService) SearchWithContext(ctx context.Context, query string, toolName string, outcome string) (*ArchetypeSearchResponse, error) {
	body := map[string]interface{}{
		"query":     query,
		"tool_name": toolName,
		"outcome":   outcome,
	}
	var result ArchetypeSearchResponse
	err := s.client.post(ctx, "/api/memory/archetype/inject", body, &result)
	return &result, err
}

// GetIndividuationScore returns the current individuation score
func (s *JungianService) GetIndividuationScore(ctx context.Context) (*IndividuationResponse, error) {
	var result IndividuationResponse
	err := s.client.get(ctx, "/api/memory/individuation", &result)
	return &result, err
}

// GetIndividuationHistory returns daily score snapshots for the specified number of days
func (s *JungianService) GetIndividuationHistory(ctx context.Context, days int) ([]IndividuationResponse, error) {
	params := url.Values{}
	params.Set("history", "true")
	params.Set("days", fmt.Sprintf("%d", days))

	path := fmt.Sprintf("/api/memory/individuation?%s", params.Encode())

	var result []IndividuationResponse
	err := s.client.get(ctx, path, &result)
	return result, err
}

// GetSynchronicityPatterns returns detected synchronicity patterns
func (s *JungianService) GetSynchronicityPatterns(ctx context.Context) ([]SynchronicityPattern, error) {
	var result []SynchronicityPattern
	err := s.client.get(ctx, "/api/memory/sync/patterns", &result)
	return result, err
}

// GetDreamHistory returns dream consolidation records
func (s *JungianService) GetDreamHistory(ctx context.Context) ([]DreamConsolidation, error) {
	var result []DreamConsolidation
	err := s.client.get(ctx, "/api/memory/dream/history", &result)
	return result, err
}
