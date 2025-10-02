package test_table


// TesttableResponse represents test_table response
type TesttableResponse struct {

	Id uint `json:"id"`

	Name string `json:"name"`

	Familyname string `json:"family_name"`

	Prefrence string `json:"prefrence"`

	Preferences string `json:"preferences"`

}

// TesttableCreateRequest represents create test_table request
type TesttableCreateRequest struct {

	Name string `json:"name" binding:""`

	Familyname string `json:"family_name" binding:""`

	Prefrence string `json:"prefrence" binding:""`

	Preferences string `json:"preferences" binding:""`

}

// TesttableUpdateRequest represents update test_table request
type TesttableUpdateRequest struct {

	Name string `json:"name" binding:"omitempty"`

	Familyname string `json:"family_name" binding:"omitempty"`

	Prefrence string `json:"prefrence" binding:"omitempty"`

	Preferences string `json:"preferences" binding:"omitempty"`

}

// PaginationResponse represents pagination response
type PaginationResponse struct {
	Data       []TesttableResponse `json:"data"`
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
