package sessions

import "time"


// SessionsResponse represents sessions response
type SessionsResponse struct {

	Id uint `json:"id"`

	Userid uint `json:"user_id"`

	Token string `json:"token"`

	Refreshtoken string `json:"refresh_token"`

	Deviceinfo string `json:"device_info"`

	Ipaddress string `json:"ip_address"`

	Useragent string `json:"user_agent"`

	Isactive bool `json:"is_active"`

	Expiresat time.Time `json:"expires_at"`

	Createdat time.Time `json:"created_at"`

	Lastusedat time.Time `json:"last_used_at"`

}

// SessionsCreateRequest represents create sessions request
type SessionsCreateRequest struct {

	Userid uint `json:"user_id" binding:"required"`

	Token string `json:"token" binding:"required,max=500"`

	Refreshtoken string `json:"refresh_token" binding:"max=500"`

	Deviceinfo string `json:"device_info" binding:""`

	Ipaddress string `json:"ip_address" binding:"max=45"`

	Useragent string `json:"user_agent" binding:""`

	Isactive bool `json:"is_active" binding:""`

	Expiresat time.Time `json:"expires_at" binding:"required"`

	Lastusedat time.Time `json:"last_used_at" binding:""`

}

// SessionsUpdateRequest represents update sessions request
type SessionsUpdateRequest struct {

	Userid uint `json:"user_id" binding:"omitempty"`

	Token string `json:"token" binding:"omitempty,max=500"`

	Refreshtoken string `json:"refresh_token" binding:"omitempty,max=500"`

	Deviceinfo string `json:"device_info" binding:"omitempty"`

	Ipaddress string `json:"ip_address" binding:"omitempty,max=45"`

	Useragent string `json:"user_agent" binding:"omitempty"`

	Isactive bool `json:"is_active" binding:"omitempty"`

	Expiresat time.Time `json:"expires_at" binding:"omitempty"`

	Lastusedat time.Time `json:"last_used_at" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []SessionsResponse `json:"data"`
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
