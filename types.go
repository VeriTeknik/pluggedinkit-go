package pluggedinkit

import "time"

// Document represents a document in the library
type Document struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Content       string            `json:"content,omitempty"`
	FileSize      int64             `json:"fileSize"`
	FileType      string            `json:"fileType"`
	Source        DocumentSource    `json:"source"`
	Tags          []string          `json:"tags,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	LastAccessAt  *time.Time        `json:"lastAccessAt,omitempty"`
	ModelProvider string            `json:"modelProvider,omitempty"`
	ModelName     string            `json:"modelName,omitempty"`
	Version       int               `json:"version"`
	Visibility    DocumentVisibility `json:"visibility"`
	Category      DocumentCategory  `json:"category,omitempty"`
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
	Source        DocumentSource `json:"source,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	Category      DocumentCategory `json:"category,omitempty"`
	DateFrom      *time.Time     `json:"dateFrom,omitempty"`
	DateTo        *time.Time     `json:"dateTo,omitempty"`
	ModelProvider string         `json:"modelProvider,omitempty"`
	ModelName     string         `json:"modelName,omitempty"`
	SearchQuery   string         `json:"searchQuery,omitempty"`
	Sort          SortOrder      `json:"sort,omitempty"`
	Limit         int            `json:"limit,omitempty"`
	Offset        int            `json:"offset,omitempty"`
}

// ListDocumentsResponse represents the response from listing documents
type ListDocumentsResponse struct {
	Documents []Document `json:"documents"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	PerPage   int        `json:"perPage"`
}

// SearchResult represents a search result
type SearchResult struct {
	DocumentID     string    `json:"documentId"`
	Title          string    `json:"title"`
	Snippet        string    `json:"snippet"`
	RelevanceScore float64   `json:"relevanceScore"`
	Tags           []string  `json:"tags,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// SearchResponse represents the response from searching documents
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Query   string         `json:"query"`
}

// CreateDocumentRequest represents a request to create a document
type CreateDocumentRequest struct {
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UpdateDocumentRequest represents a request to update a document
type UpdateDocumentRequest struct {
	Operation UpdateOperation        `json:"operation"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateDocumentResponse represents the response from updating a document
type UpdateDocumentResponse struct {
	DocumentID string    `json:"documentId"`
	Version    int       `json:"version"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// RAGResponse represents a RAG query response
type RAGResponse struct {
	Answer  string        `json:"answer"`
	Sources []RAGDocument `json:"sources,omitempty"`
}

// RAGDocument represents a document in RAG response
type RAGDocument struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Model     *ModelInfo  `json:"model,omitempty"`
	Relevance float64     `json:"relevance,omitempty"`
}

// ModelInfo represents AI model information
type ModelInfo struct {
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Version  string `json:"version,omitempty"`
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