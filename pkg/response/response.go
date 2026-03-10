package response

import (
	"github.com/gin-gonic/gin"
)

// Meta holds pagination or supplementary info
type Meta struct {
	CurrentPage       int `json:"current_page,omitempty"`
	PerPage           int `json:"per_page,omitempty"`
	Total             int `json:"total,omitempty"`
	LastPage          int `json:"last_page,omitempty"`
	CurrentPageRecord int `json:"current_page_record,omitempty"`
}

// ErrorDetail standardizes the failure format
type ErrorDetail struct {
	Type    string      `json:"type,omitempty"`    // e.g., "Validation Error", "Unauthorized"
	Message string      `json:"message,omitempty"` // General error description
	Details interface{} `json:"details,omitempty"` // Used for form validation map: map[string][]string
	Flags   interface{} `json:"flags,omitempty"`   // Legacy arbitrary flags (e.g., is_till_open, invalid_device_token)
}

// APIResponse is the unified envelope
type APIResponse struct {
	Success bool         `json:"success"`         // Explicit success/fail flag
	Code    int          `json:"code"`            // Matches HTTP status (legacy compatibility)
	Message string       `json:"message"`         // Top level human readable message
	Data    interface{}  `json:"data,omitempty"`  // Payload (Struct, Slice, Map, nil)
	Meta    *Meta        `json:"meta,omitempty"`  // Pagination data (nil if not paginated)
	Error   *ErrorDetail `json:"error,omitempty"` // Error object (nil on success)
}

// Success returns a standardized success response
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Paginated returns a standardized paginated success response
func Paginated(c *gin.Context, code int, message string, data interface{}, meta Meta) {
	c.JSON(code, APIResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// Error returns a standardized error response
func Error(c *gin.Context, code int, errType string, message string) {
	c.JSON(code, APIResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error: &ErrorDetail{
			Type:    errType,
			Message: message,
		},
	})
}

// ValidationError returns a standardized validation error response
func ValidationError(c *gin.Context, message string, details map[string][]string) {
	c.JSON(400, APIResponse{
		Success: false,
		Code:    400,
		Message: message,
		Error: &ErrorDetail{
			Type:    "Bad Request",
			Message: message,
			Details: details,
		},
	})
}

// MiddlewareError returns a standardized middleware error response with custom flags
func MiddlewareError(c *gin.Context, code int, errType string, message string, flags map[string]interface{}) {
	c.JSON(code, APIResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error: &ErrorDetail{
			Type:    errType,
			Message: message,
			Flags:   flags,
		},
	})
}
