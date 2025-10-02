package user_roles

import "time"


// UserrolesResponse represents user_roles response
type UserrolesResponse struct {

	Userid uint `json:"user_id"`

	Roleid uint `json:"role_id"`

	Assignedat time.Time `json:"assigned_at"`

	Assignedby uint `json:"assigned_by"`

}

// UserrolesCreateRequest represents create user_roles request
type UserrolesCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Roleid uint `json:"role_id" binding:"required"`

	Assignedat time.Time `json:"assigned_at" binding:""`

	Assignedby uint `json:"assigned_by" binding:""`

}

// UserrolesUpdateRequest represents update user_roles request
type UserrolesUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Roleid uint `json:"role_id" binding:"omitempty"`

	Assignedat time.Time `json:"assigned_at" binding:"omitempty"`

	Assignedby uint `json:"assigned_by" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []UserrolesResponse `json:"data"`
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
