# Plugged.in Go SDK

Official Go SDK for the Plugged.in Library API. Full support for document management, RAG queries, and file uploads.

## Installation

```bash
go get github.com/veriteknik/pluggedinkit-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    pluggedin "github.com/veriteknik/pluggedinkit-go"
)

func main() {
    // Initialize client
    client := pluggedin.NewClient("your-api-key")
    // For local development:
    // client := pluggedin.NewClientWithOptions("your-api-key", "http://localhost:12005", nil)

    ctx := context.Background()

    // List documents
    docs, err := client.Documents.List(ctx, &pluggedin.DocumentFilters{
        Source: pluggedin.SourceAll,
        Limit:  10,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d documents\n", docs.Total)

    // Search documents
    results, err := client.Documents.Search(ctx, "machine learning", nil, 10, 0)
    if err != nil {
        log.Fatal(err)
    }
    for _, result := range results.Results {
        fmt.Printf("%s - Relevance: %.2f\n", result.Title, result.RelevanceScore)
    }

    // Query RAG
    answer, err := client.RAG.AskQuestion(ctx, "What are the main features?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(answer)
}
```

## Features

- 📄 **Document Management** - Full CRUD operations for documents
- 🔍 **Semantic Search** - AI-powered document search
- 🤖 **RAG Integration** - Natural language queries to your knowledge base
- 📤 **File Uploads** - Upload files with progress tracking
- 🔄 **Version Control** - Document versioning and history
- ⚡ **Type Safety** - Full Go type definitions
- 🔐 **Authentication** - Secure API key authentication
- 🔁 **Context Support** - Context-aware operations for cancellation
- 📊 **Rate Limiting** - Built-in rate limit handling

## Authentication

Get your API key from your Plugged.in profile settings:

```go
import (
    "os"
    pluggedin "github.com/veriteknik/pluggedinkit-go"
)

client := pluggedin.NewClient(os.Getenv("PLUGGEDIN_API_KEY"))
```

## Documentation

### Document Operations

#### List Documents

```go
filters := &pluggedin.DocumentFilters{
    Source:   pluggedin.SourceAIGenerated,
    Tags:     []string{"report", "analysis"},
    Sort:     pluggedin.SortDateDesc,
    Limit:    20,
    Offset:   0,
}

response, err := client.Documents.List(ctx, filters)
if err != nil {
    log.Fatal(err)
}

for _, doc := range response.Documents {
    fmt.Printf("%s (%d bytes)\n", doc.Title, doc.FileSize)
}
```

#### Get Document

```go
// Get document metadata
doc, err := client.Documents.Get(ctx, "document-id", false, false)
if err != nil {
    log.Fatal(err)
}

// Get document with content
docWithContent, err := client.Documents.Get(ctx, "document-id", true, true)
if err != nil {
    log.Fatal(err)
}

fmt.Println(docWithContent.Content)
```

#### Search Documents

```go
filters := map[string]interface{}{
    "modelProvider": "anthropic",
    "dateFrom":      "2024-01-01T00:00:00Z",
    "tags":          "finance,q4",
}

results, err := client.Documents.Search(ctx, "quarterly report", filters, 10, 0)
if err != nil {
    log.Fatal(err)
}

for _, result := range results.Results {
    fmt.Printf("%s\n", result.Title)
    fmt.Printf("  Relevance: %.2f\n", result.RelevanceScore)
    fmt.Printf("  Snippet: %s\n", result.Snippet)
}
```

#### Update Document

```go
request := &pluggedin.UpdateDocumentRequest{
    Operation: pluggedin.OpAppend,
    Content:   "\n\n## New Section\n\nAdditional content here.",
    Metadata: map[string]interface{}{
        "changeSummary": "Added implementation details",
        "model": map[string]string{
            "name":     "claude-3-opus",
            "provider": "anthropic",
            "version":  "20240229",
        },
    },
}

response, err := client.Documents.Update(ctx, "document-id", request)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated to version %d\n", response.Version)
```

#### Create AI Document

```go
request := &pluggedin.CreateDocumentRequest{
    Title:   "API Integration Guide",
    Content: "# API Integration Guide\n\n## Introduction\n\n...",
    Metadata: map[string]interface{}{
        "format":   "md",
        "category": "documentation",
        "tags":     []string{"api", "guide"},
        "model": map[string]string{
            "name":     "gpt-4",
            "provider": "openai",
            "version":  "0613",
        },
        "prompt":     "Create an API integration guide",
        "visibility": "workspace",
    },
}

doc, err := client.Documents.Create(ctx, request)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created document: %s\n", doc.ID)
```

### RAG Operations

#### Query Knowledge Base

```go
// Simple query
answer, err := client.RAG.AskQuestion(ctx, "What are the deployment procedures?")
if err != nil {
    log.Fatal(err)
}
fmt.Println(answer)

// Query with sources
result, err := client.RAG.QueryWithSources(
    ctx,
    "Explain the authentication flow",
    "project-uuid", // Optional
)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Answer: %s\n", result.Answer)
fmt.Println("Sources:")
for _, source := range result.Sources {
    fmt.Printf("- %s (relevance: %.0f%%)\n", source.Name, source.Relevance)
}
```

#### Find Relevant Documents

```go
documents, err := client.RAG.FindRelevantDocuments(
    ctx,
    "user authentication",
    "project-uuid",
    5, // Top 5 documents
)
if err != nil {
    log.Fatal(err)
}

for _, doc := range documents {
    fmt.Printf("- %s\n", doc.Name)
    if doc.Model != nil {
        fmt.Printf("  By: %s/%s\n", doc.Model.Provider, doc.Model.Name)
    }
}
```

### File Upload Operations

#### Upload Single File

```go
metadata := &pluggedin.UploadMetadata{
    Name:        "Q4 Report.pdf",
    Description: "Quarterly financial report",
    Tags:        []string{"finance", "q4", "2024"},
    Purpose:     "Financial documentation",
    RelatedTo:   "PROJECT-123",
}

// Progress callback
onProgress := func(percent int) {
    fmt.Printf("Upload progress: %d%%\n", percent)
}

response, err := client.Uploads.UploadFile(
    ctx,
    "./report.pdf",
    metadata,
    onProgress,
)
if err != nil {
    log.Fatal(err)
}

if response.Success {
    fmt.Printf("Uploaded: %s\n", response.DocumentID)

    // Track processing
    if response.UploadID != "" {
        onUpdate := func(status *pluggedin.UploadStatus) {
            fmt.Printf("Status: %s - %s\n", status.Status, status.Message)
        }

        client.Uploads.TrackUpload(ctx, response.UploadID, onUpdate)
    }
}
```

#### Upload from Reader

```go
file, err := os.Open("document.txt")
if err != nil {
    log.Fatal(err)
}
defer file.Close()

fileInfo, _ := file.Stat()

response, err := client.Uploads.UploadReader(
    ctx,
    file,
    fileInfo.Size(),
    metadata,
    nil,
)
```

#### Batch Upload

```go
files := []pluggedin.FileUpload{
    {
        Path:     "doc1.pdf",
        Metadata: &pluggedin.UploadMetadata{Name: "doc1.pdf", Tags: []string{"batch"}},
    },
    {
        Path:     "doc2.txt",
        Metadata: &pluggedin.UploadMetadata{Name: "doc2.txt", Tags: []string{"batch"}},
    },
}

onProgress := func(current, total int) {
    fmt.Printf("Uploaded %d/%d files\n", current, total)
}

results, err := client.Uploads.UploadBatch(ctx, files, onProgress)
if err != nil {
    log.Fatal(err)
}

for i, result := range results {
    if result.Success {
        fmt.Printf("✓ %s uploaded\n", files[i].Metadata.Name)
    } else {
        fmt.Printf("✗ %s failed: %s\n", files[i].Metadata.Name, result.Error)
    }
}
```

### Error Handling

```go
doc, err := client.Documents.Get(ctx, "invalid-id", false, false)
if err != nil {
    switch e := err.(type) {
    case *pluggedin.AuthenticationError:
        fmt.Println("Invalid API key")
    case *pluggedin.RateLimitError:
        fmt.Printf("Rate limited. Retry after %d seconds\n", e.RetryAfter)
    case *pluggedin.NotFoundError:
        fmt.Printf("Document not found: %s\n", e.ID)
    case *pluggedin.ValidationError:
        fmt.Printf("Validation error: %s\n", e.Message)
    case *pluggedin.APIError:
        fmt.Printf("API error (%d): %s\n", e.StatusCode, e.Message)
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Advanced Configuration

```go
import (
    "net/http"
    "time"
)

// Custom HTTP client
httpClient := &http.Client{
    Timeout: 120 * time.Second,
    // Add custom transport, etc.
}

client := pluggedin.NewClientWithOptions(
    "your-api-key",
    "https://plugged.in",
    httpClient,
)

// Update API key at runtime
client.SetAPIKey("new-api-key")

// Change base URL
client.SetBaseURL("http://localhost:12005")
```

## Environment Variables

Store your API key securely:

```bash
export PLUGGEDIN_API_KEY=your-api-key
export PLUGGEDIN_BASE_URL=https://plugged.in
```

```go
import (
    "os"
    pluggedin "github.com/veriteknik/pluggedinkit-go"
)

client := pluggedin.NewClientWithOptions(
    os.Getenv("PLUGGEDIN_API_KEY"),
    os.Getenv("PLUGGEDIN_BASE_URL"),
    nil,
)
```

## Rate Limiting

The SDK handles rate limiting automatically:

- **API Endpoints**: 60 requests per minute
- **Document Search**: 10 requests per hour for AI document creation
- **RAG Queries**: Subject to plan limits

## Examples

See the [examples](./examples) directory for complete working examples:

- [Basic Usage](./examples/basic/main.go)
- [Document Management](./examples/documents/main.go)
- [RAG Queries](./examples/rag/main.go)
- [File Upload](./examples/upload/main.go)
- [Error Handling](./examples/errors/main.go)

## Contributing

Contributions are welcome! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

MIT - see [LICENSE](LICENSE) for details.

## Support

- **Documentation**: [https://docs.plugged.in](https://docs.plugged.in)
- **Issues**: [GitHub Issues](https://github.com/veriteknik/pluggedinkit-go/issues)
- **Discord**: [Join our community](https://discord.gg/pluggedin)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.