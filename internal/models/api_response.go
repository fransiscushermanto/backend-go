package models

// ApiResult defines the structure for a successful API response.
// 'Data' holds the actual response content (can be any Go type).
// 'Meta' is optional metadata, typically for pagination or additional context.
type ApiResult struct {
	Data interface{}             `json:"data"`           // The primary data payload of the response
	Meta *map[string]interface{} `json:"meta,omitempty"` // Optional metadata, use pointer to omit if nil
}

// ApiError defines the structure for an API error response.
// 'Message' is a human-readable summary of the error.
// 'StatusCode' is the HTTP status code (included in body for convenience).
// 'Errors' provides detailed error information for specific fields. It will be omitted if nil/empty.
// 'Data' is optional, for cases where an error response might still include some relevant data.
// 'Meta' is optional, for any additional error-related metadata.
type ApiError struct {
	StatusCode int                     `json:"status_code"` // HTTP status code
	Message    *string                 `json:"message,omitempty"`
	Errors     *map[string]string      `json:"errors,omitempty"` // Specific field errors: { "field_name": "error_message" }
	Data       *interface{}            `json:"data,omitempty"`   // Optional data field for errors, use omitempty to omit if nil
	Meta       *map[string]interface{} `json:"meta,omitempty"`   // Optional: additional error metadata, use pointer to omit if nil
}
