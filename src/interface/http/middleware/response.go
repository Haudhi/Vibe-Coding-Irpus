package middleware

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse is a generic error response structure
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render implements the render.Renderer interface
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns a 400 Bad Request error
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

// ErrBadRequest returns a 400 Bad Request error with custom message
func ErrBadRequest(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Bad request",
		ErrorText:      message,
	}
}

// ErrUnauthorized returns a 401 Unauthorized error
func ErrUnauthorized(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized",
		ErrorText:      message,
	}
}

// ErrForbidden returns a 403 Forbidden error
func ErrForbidden(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     "Forbidden",
		ErrorText:      message,
	}
}

// ErrNotFound returns a 404 Not Found error
func ErrNotFound(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Not found",
		ErrorText:      message,
	}
}

// ErrConflict returns a 409 Conflict error
func ErrConflict(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Conflict",
		ErrorText:      message,
	}
}

// ErrInternalServerError returns a 500 Internal Server Error
func ErrInternalServerError(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error",
		ErrorText:      message,
	}
}

// ErrValidation returns a 422 Unprocessable Entity error for validation failures
func ErrValidation(message string) render.Renderer {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Validation failed",
		ErrorText:      message,
	}
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Render implements the render.Renderer interface
func (s *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) render.Renderer {
	return &SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}