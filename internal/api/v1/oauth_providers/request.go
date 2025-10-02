package oauth_providers

import "time"


// OauthprovidersResponse represents oauth_providers response
type OauthprovidersResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Provider string `json:"provider"`

	Provideruserid string `json:"provider_user_id"`

	Accesstoken string `json:"access_token"`

	Refreshtoken string `json:"refresh_token"`

	Tokenexpiresat time.Time `json:"token_expires_at"`

	Profiledata string `json:"profile_data"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

}

// OauthprovidersCreateRequest represents create oauth_providers request
type OauthprovidersCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Provider string `json:"provider" binding:"required,max=50"`

	Provideruserid string `json:"provider_user_id" binding:"required,max=255"`

	Accesstoken string `json:"access_token" binding:""`

	Refreshtoken string `json:"refresh_token" binding:""`

	Tokenexpiresat time.Time `json:"token_expires_at" binding:""`

	Profiledata string `json:"profile_data" binding:""`

}

// OauthprovidersUpdateRequest represents update oauth_providers request
type OauthprovidersUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Provider string `json:"provider" binding:"omitempty,max=50"`

	Provideruserid string `json:"provider_user_id" binding:"omitempty,max=255"`

	Accesstoken string `json:"access_token" binding:"omitempty"`

	Refreshtoken string `json:"refresh_token" binding:"omitempty"`

	Tokenexpiresat time.Time `json:"token_expires_at" binding:"omitempty"`

	Profiledata string `json:"profile_data" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []OauthprovidersResponse `json:"data"`
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
