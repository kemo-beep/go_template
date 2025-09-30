package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Validation failed",
		"details": errors,
	})
}

// PaginationResponse represents paginated response data
type PaginationResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	TotalItems int64       `json:"total_items"`
}

// PaginatedResponse sends a paginated JSON response
func PaginatedResponse(c *gin.Context, data interface{}, page, limit int, totalItems int64) {
	totalPages := int(totalItems) / limit
	if int(totalItems)%limit > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, PaginationResponse{
		Success:    true,
		Data:       data,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalItems,
	})
}

// SuccessResponseData returns a success response (for use with c.JSON)
func SuccessResponseData(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponseData returns an error response (for use with c.JSON)
func ErrorResponseData(message string) Response {
	return Response{
		Success: false,
		Error:   message,
	}
}
