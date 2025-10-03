package wow_table


// WowtableResponse represents wow_table response
type WowtableResponse struct {

	Id uint `json:"id"`

	Swim string `json:"swim"`

}

// WowtableCreateRequest represents create wow_table request
type WowtableCreateRequest struct {

	Swim string `json:"swim" binding:""`

}

// WowtableUpdateRequest represents update wow_table request
type WowtableUpdateRequest struct {

	Swim string `json:"swim" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []WowtableResponse `json:"data"`
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
