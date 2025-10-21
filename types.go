package pluggedinkit

import "time"

// Document represents a document in the library
type Document struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description,omitempty"`
	FileName          string             `json:"fileName"`
	FileSize          int64              `json:"fileSize"`
	MimeType          string             `json:"mimeType"`
	Tags              []string           `json:"tags,omitempty"`
	Source            DocumentSource     `json:"source"`
	Visibility        DocumentVisibility `json:"visibility"`
	Version           int                `json:"version"`
	CreatedAt         time.Time          `json:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt"`
	AIMetadata        map[string]any     `json:"aiMetadata,omitempty"`
	ModelAttributions []ModelAttribution `json:"modelAttributions,omitempty"`
	Content           string             `json:"content,omitempty"`
	ContentEncoding   string             `json:"contentEncoding,omitempty"`
	ContentHash       string             `json:"contentHash,omitempty"`
	ParentDocumentID  string             `json:"parentDocumentId,omitempty"`
}

// ModelAttribution represents model contribution metadata
type ModelAttribution struct {
	ModelName        string         `json:"modelName"`
	ModelProvider    string         `json:"modelProvider"`
	ContributionType string         `json:"contributionType"`
	Timestamp        time.Time      `json:"timestamp"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

// DocumentSource represents the source of a document
type DocumentSource string

const (
	SourceAll         DocumentSource = "all"
	SourceUpload      DocumentSource = "upload"
	SourceAIGenerated DocumentSource = "ai_generated"
	SourceAPI         DocumentSource = "api"
)

// DocumentVisibility represents document visibility
type DocumentVisibility string

const (
	VisibilityPrivate   DocumentVisibility = "private"
	VisibilityWorkspace DocumentVisibility = "workspace"
	VisibilityPublic    DocumentVisibility = "public"
)

// DocumentCategory represents document category
type DocumentCategory string

const (
	CategoryReport        DocumentCategory = "report"
	CategoryAnalysis      DocumentCategory = "analysis"
	CategoryDocumentation DocumentCategory = "documentation"
	CategoryGuide         DocumentCategory = "guide"
	CategoryResearch      DocumentCategory = "research"
	CategoryCode          DocumentCategory = "code"
	CategoryOther         DocumentCategory = "other"
)

// DocumentFormat represents document format
type DocumentFormat string

const (
	FormatMarkdown DocumentFormat = "md"
	FormatText     DocumentFormat = "txt"
	FormatJSON     DocumentFormat = "json"
	FormatHTML     DocumentFormat = "html"
)

// SortOrder represents sort order for listing
type SortOrder string

const (
	SortDateDesc SortOrder = "date_desc"
	SortDateAsc  SortOrder = "date_asc"
	SortTitle    SortOrder = "title"
	SortSize     SortOrder = "size"
)

// UpdateOperation represents update operations
type UpdateOperation string

const (
	OpReplace UpdateOperation = "replace"
	OpAppend  UpdateOperation = "append"
	OpPrepend UpdateOperation = "prepend"
)

// DocumentFilters represents filters for listing documents
type DocumentFilters struct {
	Source        DocumentSource   `json:"source,omitempty"`
	Tags          []string         `json:"tags,omitempty"`
	Category      DocumentCategory `json:"category,omitempty"`
	DateFrom      *time.Time       `json:"dateFrom,omitempty"`
	DateTo        *time.Time       `json:"dateTo,omitempty"`
	ModelProvider string           `json:"modelProvider,omitempty"`
	ModelName     string           `json:"modelName,omitempty"`
	SearchQuery   string           `json:"searchQuery,omitempty"`
	Sort          SortOrder        `json:"sort,omitempty"`
	Limit         int              `json:"limit,omitempty"`
	Offset        int              `json:"offset,omitempty"`
}

// ListDocumentsResponse represents the response from listing documents
type ListDocumentsResponse struct {
	Documents []Document `json:"documents"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

// SearchResult represents a search result
type SearchResult struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description,omitempty"`
	Snippet           string             `json:"snippet"`
	RelevanceScore    float64            `json:"relevanceScore"`
	Source            string             `json:"source"`
	AIMetadata        map[string]any     `json:"aiMetadata,omitempty"`
	Tags              []string           `json:"tags,omitempty"`
	Visibility        string             `json:"visibility"`
	CreatedAt         time.Time          `json:"createdAt"`
	ModelAttributions []ModelAttribution `json:"modelAttributions,omitempty"`
}

// SearchResponse represents the response from searching documents
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
	HasMore bool           `json:"hasMore"`
}

// CreateDocumentRequest represents a request to create a document
type CreateDocumentRequest struct {
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Format   string                 `json:"format,omitempty"`
	Category string                 `json:"category,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateDocumentRequest represents a request to update a document
type UpdateDocumentRequest struct {
	Operation UpdateOperation        `json:"operation"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateDocumentResponse represents the response from updating a document
type UpdateDocumentResponse struct {
	Success     bool   `json:"success"`
	DocumentID  string `json:"documentId"`
	Version     int    `json:"version"`
	FileWritten bool   `json:"fileWritten"`
	Message     string `json:"message,omitempty"`
}

// RAGResponse represents a RAG query response
type RAGResponse struct {
	Success     bool     `json:"success"`
	Answer      string   `json:"answer,omitempty"`
	Sources     []string `json:"sources,omitempty"`
	DocumentIDs []string `json:"documentIds,omitempty"`
	Error       string   `json:"error,omitempty"`
}

// RAGDocumentReference represents a matched document reference
type RAGDocumentReference struct {
	DocumentID string `json:"documentId"`
	Source     string `json:"source,omitempty"`
}

// RAGStorageStats represents storage metrics returned by the API
type RAGStorageStats struct {
	DocumentsCount     int     `json:"documents_count"`
	TotalChunks        int     `json:"total_chunks"`
	EstimatedStorageMB float64 `json:"estimated_storage_mb"`
	VectorsCount       int     `json:"vectors_count,omitempty"`
	EmbeddingDimension int     `json:"embedding_dimension,omitempty"`
	IsEstimate         bool    `json:"is_estimate,omitempty"`
}

// UploadMetadata represents metadata for file uploads
type UploadMetadata struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Purpose     string   `json:"purpose,omitempty"`
	RelatedTo   string   `json:"relatedTo,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
}

// UploadResponse represents the response from uploading a file
type UploadResponse struct {
	Success    bool   `json:"success"`
	DocumentID string `json:"documentId,omitempty"`
	UploadID   string `json:"uploadId,omitempty"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message,omitempty"`
}

// UploadStatus represents the status of an upload
type UploadStatus struct {
	UploadID string `json:"uploadId"`
	Status   string `json:"status"`
	Message  string `json:"message"`
	Progress int    `json:"progress"`
}
