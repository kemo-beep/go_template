package permissions

import "time"


// PermissionsResponse represents permissions response
type PermissionsResponse struct {

	Id uint `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	Resource string `json:"resource"`

	Action string `json:"action"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

}

// PermissionsCreateRequest represents create permissions request
type PermissionsCreateRequest struct {

	Name string `json:"name" binding:"required,max=100"`

	Description string `json:"description" binding:""`

	Resource string `json:"resource" binding:"required,max=100"`

	Action string `json:"action" binding:"required,max=50"`

}

// PermissionsUpdateRequest represents update permissions request
type PermissionsUpdateRequest struct {

	Name string `json:"name" binding:"omitempty,max=100"`

	Description string `json:"description" binding:"omitempty"`

	Resource string `json:"resource" binding:"omitempty,max=100"`

	Action string `json:"action" binding:"omitempty,max=50"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []PermissionsResponse `json:"data"`
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
