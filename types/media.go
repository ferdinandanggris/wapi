package types

// MediaUploadResponse contains the uploaded media ID.
type MediaUploadResponse struct {
	ID string `json:"id"`
}

// MediaInfo contains metadata about an uploaded media file.
type MediaInfo struct {
	URL              string `json:"url"`
	MimeType         string `json:"mime_type"`
	SHA256           string `json:"sha256"`
	FileSize         int64  `json:"file_size"`
	ID               string `json:"id"`
	MessagingProduct string `json:"messaging_product"`
}
