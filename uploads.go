package pluggedinkit

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
)

// UploadsService handles upload operations.
// The public API no longer exposes direct binary upload endpoints. The methods
// now return informative errors to help consumers migrate to the supported
// flows.
type UploadsService struct {
	client *Client
}

// ProgressCallback is retained for backwards compatibility.
type ProgressCallback func(percent int)

// UploadFile is no longer supported.
func (s *UploadsService) UploadFile(ctx context.Context, filePath string, metadata *UploadMetadata, onProgress ProgressCallback) (*UploadResponse, error) {
	_ = ctx
	_ = metadata
	_ = onProgress

	fileName := filepath.Base(filePath)
	return nil, fmt.Errorf("binary upload for %s is no longer supported via the API", fileName)
}

// UploadReader is no longer supported.
func (s *UploadsService) UploadReader(ctx context.Context, reader io.Reader, size int64, metadata *UploadMetadata, onProgress ProgressCallback) (*UploadResponse, error) {
	_ = ctx
	_ = reader
	_ = size
	_ = metadata
	_ = onProgress

	return nil, fmt.Errorf("binary uploads are no longer supported via the API")
}

// UploadBatch is no longer supported.
func (s *UploadsService) UploadBatch(ctx context.Context, files []FileUpload, onProgress func(current, total int)) ([]UploadResponse, error) {
	_ = ctx
	_ = files
	_ = onProgress

	return nil, fmt.Errorf("batch uploads are no longer supported via the API")
}

// TrackUpload is no longer supported.
func (s *UploadsService) TrackUpload(ctx context.Context, uploadID string, onUpdate func(status *UploadStatus)) error {
	_ = ctx
	_ = uploadID
	_ = onUpdate

	return fmt.Errorf("upload status tracking is no longer available via the API")
}

// FileUpload represents a file to upload (retained for compatibility).
type FileUpload struct {
	Path     string
	Metadata *UploadMetadata
}
