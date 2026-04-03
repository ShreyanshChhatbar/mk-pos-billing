package response

import (
	"encoding/json"

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
	Success bool                   `json:"success"`         // Explicit success/fail flag
	Code    int                    `json:"code"`            // Matches HTTP status (legacy compatibility)
	Type    string                 `json:"type,omitempty"`    // e.g., "Validation Error", "Unauthorized"
	Message string                 `json:"message"`         // Top level human readable message
	Data    interface{}            `json:"data,omitempty"`  // Payload (Struct, Slice, Map, nil)
	Meta    *Meta                  `json:"meta,omitempty"`  // Pagination data (nil if not paginated)
	Error   *ErrorDetail           `json:"error,omitempty"` // Error object (nil on success)
	Flags   map[string]interface{} `json:"-"`               // Custom top-level flags (flattened via MarshalJSON)
}

// MarshalJSON flattens the Flags map into the top-level JSON object
func (r APIResponse) MarshalJSON() ([]byte, error) {
	type Alias APIResponse
	base, err := json.Marshal(Alias(r))
	if err != nil {
		return nil, err
	}

	if len(r.Flags) == 0 {
		return base, nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}

	for k, v := range r.Flags {
		m[k] = v
	}

	return json.Marshal(m)
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

// Error returns a standardized error response with optional top-level flags
func Error(c *gin.Context, code int, message string, flags ...map[string]interface{}) {
	res := APIResponse{
		Code:    code,
		Type:    resolveResponseType(code),
		Message: message,
		Success: false,
	}

	if len(flags) > 0 {
		res.Flags = flags[0]
	}

	c.JSON(code, res)
}

// ValidationError returns a standardized validation error response
func ValidationError(c *gin.Context, message string, details map[string][]string) {
	c.JSON(400, APIResponse{
		Success: false,
		Code:    400,
		Message: message,
		Error: &ErrorDetail{
			Type:    resolveResponseType(400),
			Message: message,
			Details: details,
		},
	})
}

// MiddlewareError returns a standardized middleware error response with custom top-level flags
func MiddlewareError(c *gin.Context, code int, errType string, message string, flags map[string]interface{}) {
	c.JSON(code, APIResponse{
		Success: false,
		Code:    code,
		Message: message,
		Type:    resolveResponseType(code),
		Flags:   flags,
	})
}

func resolveResponseType(code int) string {
	mapCodes := map[int]string{
		200: "Success",
		201: "Created",
		202: "Success",
		307: "Temporary Redirect",
		400: "Bad request",
		401: "Unauthorized (Invalid token)",
		403: "Forbidden",
		404: "Url not found",
		405: "Method not allowed",
		406: "Not acceptable",
		408: "Invalid Device ID",
		409: "Order not acceptable",
		429: "Too many requests",
		500: "Server error",
	}

	if val, ok := mapCodes[code]; ok {
		return val
	}
	return "Internal Server Error"
}