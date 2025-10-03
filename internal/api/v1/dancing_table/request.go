package dancing_table


// DancingtableResponse represents dancing_table response
type DancingtableResponse struct {

	Id uint `json:"id"`

	Frequency string `json:"frequency"`

}

// DancingtableCreateRequest represents create dancing_table request
type DancingtableCreateRequest struct {

	Frequency string `json:"frequency" binding:""`

}

// DancingtableUpdateRequest represents update dancing_table request
type DancingtableUpdateRequest struct {

	Frequency string `json:"frequency" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []DancingtableResponse `json:"data"`
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
