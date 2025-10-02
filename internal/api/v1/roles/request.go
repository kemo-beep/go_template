package roles

import "time"


// RolesResponse represents roles response
type RolesResponse struct {

	Id uint `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

}

// RolesCreateRequest represents create roles request
type RolesCreateRequest struct {

	Name string `json:"name" binding:"required,max=50"`

	Description string `json:"description" binding:""`

}

// RolesUpdateRequest represents update roles request
type RolesUpdateRequest struct {

	Name string `json:"name" binding:"omitempty,max=50"`

	Description string `json:"description" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []RolesResponse `json:"data"`
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
