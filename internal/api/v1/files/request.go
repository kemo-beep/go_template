package files

import "time"


// FilesResponse represents files response
type FilesResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Filename string `json:"file_name"`

	Filesize uint `json:"file_size"`

	Filetype string `json:"file_type"`

	R2key string `json:"r2_key"`

	R2url string `json:"r2_url"`

	Ispublic bool `json:"is_public"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

	Deletedat time.Time `json:"deleted_at"`

}

// FilesCreateRequest represents create files request
type FilesCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Filename string `json:"file_name" binding:"required"`

	Filesize uint `json:"file_size" binding:"required"`

	Filetype string `json:"file_type" binding:"required"`

	R2key string `json:"r2_key" binding:"required"`

	R2url string `json:"r2_url" binding:"required"`

	Ispublic bool `json:"is_public" binding:""`

	Deletedat time.Time `json:"deleted_at" binding:""`

}

// FilesUpdateRequest represents update files request
type FilesUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Filename string `json:"file_name" binding:"omitempty"`

	Filesize uint `json:"file_size" binding:"omitempty"`

	Filetype string `json:"file_type" binding:"omitempty"`

	R2key string `json:"r2_key" binding:"omitempty"`

	R2url string `json:"r2_url" binding:"omitempty"`

	Ispublic bool `json:"is_public" binding:"omitempty"`

	Deletedat time.Time `json:"deleted_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []FilesResponse `json:"data"`
	Pagination PaginationInfo            `json:"pagination"`
}

// PaginationInfo represents pagination information
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}
