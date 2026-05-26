package cloud

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/ferdinandanggris/wapi/types"
)

// UploadMedia uploads a file using multipart/form-data.
// Supported types: image/png, image/jpeg, image/webp, video/mp4, audio/ogg, audio/mp3, application/pdf.
func (c *CloudClient) UploadMedia(ctx context.Context, phoneNumberID, filename string, file io.Reader, mimeType string) (*types.MediaUploadResponse, error) {
	path := fmt.Sprintf("%s/media", phoneNumberID)
	var resp types.MediaUploadResponse
	if err := c.doUpload(ctx, path, filename, mimeType, file, &resp); err != nil {
		return nil, fmt.Errorf("upload media: %w", err)
	}
	return &resp, nil
}

// GetMediaURL returns metadata (including download URL) for a media ID.
func (c *CloudClient) GetMediaURL(ctx context.Context, mediaID string) (*types.MediaInfo, error) {
	var info types.MediaInfo
	if err := c.do(ctx, "GET", mediaID, nil, &info); err != nil {
		return nil, fmt.Errorf("get media url: %w", err)
	}
	return &info, nil
}

// DownloadMedia downloads a media file by its ID. Caller must close the returned body.
func (c *CloudClient) DownloadMedia(ctx context.Context, mediaID string) (io.ReadCloser, error) {
	info, err := c.GetMediaURL(ctx, mediaID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", info.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("download media: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download media: %w", err)
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		return nil, fmt.Errorf("download media: HTTP %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// DeleteMedia deletes an uploaded media file by its ID.
func (c *CloudClient) DeleteMedia(ctx context.Context, mediaID string) error {
	if err := c.doDelete(ctx, mediaID); err != nil {
		return fmt.Errorf("delete media: %w", err)
	}
	return nil
}
