package role_permissions

import "time"


// RolepermissionsResponse represents role_permissions response
type RolepermissionsResponse struct {

	Roleid uint `json:"role_id"`

	Permissionid uint `json:"permission_id"`

	Createdat time.Time `json:"created_at"`

}

// RolepermissionsCreateRequest represents create role_permissions request
type RolepermissionsCreateRequest struct {

	Roleid uint `json:"role_id" binding:"required"`

	Permissionid uint `json:"permission_id" binding:"required"`

}

// RolepermissionsUpdateRequest represents update role_permissions request
type RolepermissionsUpdateRequest struct {

	Roleid uint `json:"role_id" binding:"omitempty"`

	Permissionid uint `json:"permission_id" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []RolepermissionsResponse `json:"data"`
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
