package refresh_tokens

import "time"


// RefreshtokensResponse represents refresh_tokens response
type RefreshtokensResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Token string `json:"token"`

	Expiresat time.Time `json:"expires_at"`

	Isrevoked bool `json:"is_revoked"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

	Deletedat time.Time `json:"deleted_at"`

}

// RefreshtokensCreateRequest represents create refresh_tokens request
type RefreshtokensCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Token string `json:"token" binding:"required"`

	Expiresat time.Time `json:"expires_at" binding:"required"`

	Isrevoked bool `json:"is_revoked" binding:""`

	Deletedat time.Time `json:"deleted_at" binding:""`

}

// RefreshtokensUpdateRequest represents update refresh_tokens request
type RefreshtokensUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Token string `json:"token" binding:"omitempty"`

	Expiresat time.Time `json:"expires_at" binding:"omitempty"`

	Isrevoked bool `json:"is_revoked" binding:"omitempty"`

	Deletedat time.Time `json:"deleted_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []RefreshtokensResponse `json:"data"`
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
