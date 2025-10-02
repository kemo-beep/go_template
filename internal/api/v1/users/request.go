package users

import "time"


// UsersResponse represents users response
type UsersResponse struct {

	Id uint `json:"id"`

	Email string `json:"email"`

	Password string `json:"password"`

	Name string `json:"name"`

	Isactive bool `json:"is_active"`

	Isadmin bool `json:"is_admin"`

	Createdat time.Time `json:"created_at"`

	Updatedat time.Time `json:"updated_at"`

	Deletedat time.Time `json:"deleted_at"`

	Emailverified bool `json:"email_verified"`

	Emailverifiedat time.Time `json:"email_verified_at"`

	Lastloginat time.Time `json:"last_login_at"`

	Failedloginattempts uint `json:"failed_login_attempts"`

	Lockeduntil time.Time `json:"locked_until"`

	Metadata string `json:"metadata"`

	Nickname string `json:"nickname"`

	Bio string `json:"bio"`

}

// UsersCreateRequest represents create users request
type UsersCreateRequest struct {

	Email string `json:"email" binding:"required"`

	Password string `json:"password" binding:"required"`

	Name string `json:"name" binding:"required"`

	Isactive bool `json:"is_active" binding:""`

	Isadmin bool `json:"is_admin" binding:""`

	Deletedat time.Time `json:"deleted_at" binding:""`

	Emailverified bool `json:"email_verified" binding:""`

	Emailverifiedat time.Time `json:"email_verified_at" binding:""`

	Lastloginat time.Time `json:"last_login_at" binding:""`

	Failedloginattempts uint `json:"failed_login_attempts" binding:""`

	Lockeduntil time.Time `json:"locked_until" binding:""`

	Metadata string `json:"metadata" binding:""`

	Nickname string `json:"nickname" binding:"max=50"`

	Bio string `json:"bio" binding:""`

}

// UsersUpdateRequest represents update users request
type UsersUpdateRequest struct {

	Email string `json:"email" binding:"omitempty"`

	Password string `json:"password" binding:"omitempty"`

	Name string `json:"name" binding:"omitempty"`

	Isactive bool `json:"is_active" binding:"omitempty"`

	Isadmin bool `json:"is_admin" binding:"omitempty"`

	Deletedat time.Time `json:"deleted_at" binding:"omitempty"`

	Emailverified bool `json:"email_verified" binding:"omitempty"`

	Emailverifiedat time.Time `json:"email_verified_at" binding:"omitempty"`

	Lastloginat time.Time `json:"last_login_at" binding:"omitempty"`

	Failedloginattempts uint `json:"failed_login_attempts" binding:"omitempty"`

	Lockeduntil time.Time `json:"locked_until" binding:"omitempty"`

	Metadata string `json:"metadata" binding:"omitempty"`

	Nickname string `json:"nickname" binding:"omitempty,max=50"`

	Bio string `json:"bio" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []UsersResponse `json:"data"`
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
