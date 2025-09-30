package files

// FileResponse represents file response
type FileResponse struct {
	ID        uint   `json:"id"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	FileType  string `json:"file_type"`
	R2URL     string `json:"r2_url"`
	IsPublic  bool   `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

// UploadResponse represents upload response
type UploadResponse struct {
	File        FileResponse `json:"file"`
	DownloadURL string       `json:"download_url,omitempty"`
}

// DownloadURLResponse represents download URL response
type DownloadURLResponse struct {
	URL       string `json:"url"`
	ExpiresIn int    `json:"expires_in"` // seconds
}
