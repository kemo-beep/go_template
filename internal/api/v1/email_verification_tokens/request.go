package email_verification_tokens

import "time"


// EmailverificationtokensResponse represents email_verification_tokens response
type EmailverificationtokensResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Email string `json:"email"`

	Token string `json:"token"`

	Expiresat time.Time `json:"expires_at"`

	Used bool `json:"used"`

	Usedat time.Time `json:"used_at"`

	Createdat time.Time `json:"created_at"`

}

// EmailverificationtokensCreateRequest represents create email_verification_tokens request
type EmailverificationtokensCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Email string `json:"email" binding:"required,max=255"`

	Token string `json:"token" binding:"required,max=255"`

	Expiresat time.Time `json:"expires_at" binding:"required"`

	Used bool `json:"used" binding:""`

	Usedat time.Time `json:"used_at" binding:""`

}

// EmailverificationtokensUpdateRequest represents update email_verification_tokens request
type EmailverificationtokensUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Email string `json:"email" binding:"omitempty,max=255"`

	Token string `json:"token" binding:"omitempty,max=255"`

	Expiresat time.Time `json:"expires_at" binding:"omitempty"`

	Used bool `json:"used" binding:"omitempty"`

	Usedat time.Time `json:"used_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []EmailverificationtokensResponse `json:"data"`
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
