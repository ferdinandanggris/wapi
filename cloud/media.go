package cloud

import (
	"bytes"
	"context"
	"encoding/json"
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

// ResumableUpload uses Meta's Resumable Upload API to upload a file and return a handle.
// Two steps: (1) POST /{appId}/uploads to get session ID, (2) POST /{sessionId} with file content.
// The returned handle is used as profile_picture_handle in UpdateBusinessProfile.
func (c *CloudClient) ResumableUpload(ctx context.Context, appID string, data []byte, mimeType string) (string, error) {
	// Step 1: Initialize upload session
	initPath := fmt.Sprintf("%s/uploads?file_length=%d&file_type=%s", appID, len(data), mimeType)
	var initResp struct {
		ID string `json:"id"`
	}
	if err := c.do(ctx, "POST", initPath, nil, &initResp); err != nil {
		return "", fmt.Errorf("init resumable upload: %w", err)
	}
	if initResp.ID == "" {
		return "", fmt.Errorf("init resumable upload: empty session id")
	}

	// Step 2: Upload file content to the session endpoint
	uploadReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL(initResp.ID), bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("create upload request: %w", err)
	}
	uploadReq.Header.Set("Authorization", "Bearer "+c.accessToken)
	uploadReq.Header.Set("file_offset", "0")
	uploadReq.Header.Set("Content-Type", "application/octet-stream")
	uploadReq.ContentLength = int64(len(data))

	var uploadResp struct {
		Handle string `json:"h"`
	}
	// Use sendRequest directly since we bypass the standard do() for the session ID endpoint
	rawResp, err := c.httpClient.Do(uploadReq)
	if err != nil {
		return "", fmt.Errorf("send upload request: %w", err)
	}
	defer rawResp.Body.Close()

	body, err := io.ReadAll(rawResp.Body)
	if err != nil {
		return "", fmt.Errorf("read upload response: %w", err)
	}
	if rawResp.StatusCode >= 400 {
		return "", fmt.Errorf("upload content: HTTP %d: %s", rawResp.StatusCode, string(body))
	}
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return "", fmt.Errorf("parse upload response: %w", err)
	}
	if uploadResp.Handle == "" {
		return "", fmt.Errorf("upload content: empty handle in response: %s", string(body))
	}

	return uploadResp.Handle, nil
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
