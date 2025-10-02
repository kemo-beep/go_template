package user_2fa

import "time"


// User2faResponse represents user_2fa response
type User2faResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Secret string `json:"secret"`

	Backupcodes string `json:"backup_codes"`

	Isenabled bool `json:"is_enabled"`

	Enabledat time.Time `json:"enabled_at"`

	Lastusedat time.Time `json:"last_used_at"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

}

// User2faCreateRequest represents create user_2fa request
type User2faCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Secret string `json:"secret" binding:"required,max=255"`

	Backupcodes string `json:"backup_codes" binding:""`

	Isenabled bool `json:"is_enabled" binding:""`

	Enabledat time.Time `json:"enabled_at" binding:""`

	Lastusedat time.Time `json:"last_used_at" binding:""`

}

// User2faUpdateRequest represents update user_2fa request
type User2faUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Secret string `json:"secret" binding:"omitempty,max=255"`

	Backupcodes string `json:"backup_codes" binding:"omitempty"`

	Isenabled bool `json:"is_enabled" binding:"omitempty"`

	Enabledat time.Time `json:"enabled_at" binding:"omitempty"`

	Lastusedat time.Time `json:"last_used_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []User2faResponse `json:"data"`
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
