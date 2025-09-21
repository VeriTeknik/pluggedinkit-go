package pluggedinkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadsService handles file upload operations
type UploadsService struct {
	client *Client
}

// ProgressCallback is called during upload progress
type ProgressCallback func(percent int)

// UploadFile uploads a file from a file path
func (s *UploadsService) UploadFile(ctx context.Context, filePath string, metadata *UploadMetadata, onProgress ProgressCallback) (*UploadResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = &UploadMetadata{}
	}
	if metadata.Name == "" {
		metadata.Name = filepath.Base(filePath)
	}

	return s.uploadReader(ctx, file, fileInfo.Size(), metadata, onProgress)
}

// UploadReader uploads a file from an io.Reader
func (s *UploadsService) UploadReader(ctx context.Context, reader io.Reader, size int64, metadata *UploadMetadata, onProgress ProgressCallback) (*UploadResponse, error) {
	return s.uploadReader(ctx, reader, size, metadata, onProgress)
}

// uploadReader performs the actual upload
func (s *UploadsService) uploadReader(ctx context.Context, reader io.Reader, size int64, metadata *UploadMetadata, onProgress ProgressCallback) (*UploadResponse, error) {
	// Create multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file field
	part, err := writer.CreateFormFile("file", metadata.Name)
	if err != nil {
		return nil, err
	}

	// Track progress if callback provided
	var written int64
	if onProgress != nil {
		reader = &progressReader{
			reader:   reader,
			size:     size,
			callback: onProgress,
			written:  &written,
		}
	}

	// Copy file content
	if _, err := io.Copy(part, reader); err != nil {
		return nil, err
	}

	// Add metadata field
	if metadata != nil {
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return nil, err
		}
		if err := writer.WriteField("metadata", string(metadataJSON)); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/library/upload", s.client.baseURL), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var uploadResp UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return &uploadResp, &APIError{
			StatusCode: resp.StatusCode,
			Message:    uploadResp.Error,
		}
	}

	return &uploadResp, nil
}

// UploadBatch uploads multiple files
func (s *UploadsService) UploadBatch(ctx context.Context, files []FileUpload, onProgress func(current, total int)) ([]UploadResponse, error) {
	results := make([]UploadResponse, len(files))

	for i, file := range files {
		if onProgress != nil {
			onProgress(i, len(files))
		}

		result, err := s.UploadFile(ctx, file.Path, file.Metadata, nil)
		if err != nil {
			results[i] = UploadResponse{
				Success: false,
				Error:   err.Error(),
			}
		} else {
			results[i] = *result
		}
	}

	if onProgress != nil {
		onProgress(len(files), len(files))
	}

	return results, nil
}

// TrackUpload tracks the status of an upload
func (s *UploadsService) TrackUpload(ctx context.Context, uploadID string, onUpdate func(status *UploadStatus)) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var status UploadStatus
			path := fmt.Sprintf("/api/library/upload/%s/status", uploadID)
			if err := s.client.get(ctx, path, &status); err != nil {
				return err
			}

			if onUpdate != nil {
				onUpdate(&status)
			}

			if status.Status == "completed" || status.Status == "failed" {
				return nil
			}
		}
	}
}

// FileUpload represents a file to upload
type FileUpload struct {
	Path     string
	Metadata *UploadMetadata
}

// progressReader wraps an io.Reader to track progress
type progressReader struct {
	reader   io.Reader
	size     int64
	written  *int64
	callback ProgressCallback
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		*pr.written += int64(n)
		if pr.callback != nil && pr.size > 0 {
			percent := int((*pr.written * 100) / pr.size)
			pr.callback(percent)
		}
	}
	return n, err
}