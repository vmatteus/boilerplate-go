package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Message    string      `json:"message,omitempty"`
}

// Pagination represents pagination info
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// JSON sends a JSON response
func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
}

// Success sends a success response
func Success(c *gin.Context, data interface{}, message ...string) {
	response := SuccessResponse{
		Data: data,
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(http.StatusOK, response)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, err string, message ...string) {
	response := ErrorResponse{
		Error: err,
	}
	if len(message) > 0 {
		response.Message = message[0]
	}
	c.JSON(statusCode, response)
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, err string, message ...string) {
	Error(c, http.StatusBadRequest, err, message...)
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, err string, message ...string) {
	Error(c, http.StatusNotFound, err, message...)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, err string, message ...string) {
	Error(c, http.StatusInternalServerError, err, message...)
}

// Paginated sends a paginated response
func Paginated(c *gin.Context, data interface{}, page, limit int, total int64, message ...string) {
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
	if len(message) > 0 {
		response.Message = message[0]
	}

	c.JSON(http.StatusOK, response)
}
