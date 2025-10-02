package api_keys

import "time"


// ApikeysResponse represents api_keys response
type ApikeysResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Name string `json:"name"`

	Keyhash string `json:"key_hash"`

	Prefix string `json:"prefix"`

	Scopes string `json:"scopes"`

	Ratelimit uint `json:"rate_limit"`

	Isactive bool `json:"is_active"`

	Lastusedat time.Time `json:"last_used_at"`

	Expiresat time.Time `json:"expires_at"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

}

// ApikeysCreateRequest represents create api_keys request
type ApikeysCreateRequest struct {

	Userid uint `json:"user_id" binding:""`

	Name string `json:"name" binding:"required,max=100"`

	Keyhash string `json:"key_hash" binding:"required,max=255"`

	Prefix string `json:"prefix" binding:"required,max=20"`

	Scopes string `json:"scopes" binding:""`

	Ratelimit uint `json:"rate_limit" binding:""`

	Isactive bool `json:"is_active" binding:""`

	Lastusedat time.Time `json:"last_used_at" binding:""`

	Expiresat time.Time `json:"expires_at" binding:""`

}

// ApikeysUpdateRequest represents update api_keys request
type ApikeysUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Name string `json:"name" binding:"omitempty,max=100"`

	Keyhash string `json:"key_hash" binding:"omitempty,max=255"`

	Prefix string `json:"prefix" binding:"omitempty,max=20"`

	Scopes string `json:"scopes" binding:"omitempty"`

	Ratelimit uint `json:"rate_limit" binding:"omitempty"`

	Isactive bool `json:"is_active" binding:"omitempty"`

	Lastusedat time.Time `json:"last_used_at" binding:"omitempty"`

	Expiresat time.Time `json:"expires_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []ApikeysResponse `json:"data"`
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
