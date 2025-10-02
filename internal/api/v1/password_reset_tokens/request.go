package password_reset_tokens

import "time"


// PasswordresettokensResponse represents password_reset_tokens response
type PasswordresettokensResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Token string `json:"token"`

	Expiresat time.Time `json:"expires_at"`

	Used bool `json:"used"`

	Usedat time.Time `json:"used_at"`

	Createdat time.Time `json:"created_at"`

}

// PasswordresettokensCreateRequest represents create password_reset_tokens request
type PasswordresettokensCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Token string `json:"token" binding:"required,max=255"`

	Expiresat time.Time `json:"expires_at" binding:"required"`

	Used bool `json:"used" binding:""`

	Usedat time.Time `json:"used_at" binding:""`

}

// PasswordresettokensUpdateRequest represents update password_reset_tokens request
type PasswordresettokensUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Token string `json:"token" binding:"omitempty,max=255"`

	Expiresat time.Time `json:"expires_at" binding:"omitempty"`

	Used bool `json:"used" binding:"omitempty"`

	Usedat time.Time `json:"used_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []PasswordresettokensResponse `json:"data"`
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
