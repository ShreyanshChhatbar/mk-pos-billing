package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

// RespondJSON sends a standardized success response with optional meta block (reserve meta for pagination).
func RespondJSON(c *gin.Context, code int, message string, data interface{}, meta interface{}) {
	c.JSON(code, APIResponse{
		Success:    true,
		StatusCode: code,
		Message:    message,
		Data:       data,
		Meta:       meta,
	})
}

// RespondError sends a standardized error response.
func RespondError(c *gin.Context, code int, message string, errors interface{}) {
	c.JSON(code, APIResponse{
		Success:    false,
		StatusCode: code,
		Message:    message,
		Errors:     errors,
	})
}
