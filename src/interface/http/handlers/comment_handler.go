package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/application/usecases"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/interface/http/middleware"
)

// CommentHandler handles HTTP requests for comment operations
type CommentHandler struct {
	addCommentUseCase *usecases.AddCommentUseCase
	getCommentsUseCase *usecases.GetCommentsUseCase
}

// NewCommentHandler creates a new comment handler
func NewCommentHandler(
	addCommentUseCase *usecases.AddCommentUseCase,
	getCommentsUseCase *usecases.GetCommentsUseCase,
) *CommentHandler {
	return &CommentHandler{
		addCommentUseCase:  addCommentUseCase,
		getCommentsUseCase: getCommentsUseCase,
	}
}

// AddCommentRequest represents the request body for adding a comment
type AddCommentRequest struct {
	Content string `json:"content" validate:"required,max=2000"`
}

// AddComment handles adding a comment to a ticket
func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketID")
	if ticketID == "" {
		render.Render(w, r, middleware.ErrBadRequest("ticket ID is required"))
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(ticketID); err != nil {
		render.Render(w, r, middleware.ErrBadRequest("invalid ticket ID format"))
		return
	}

	var req AddCommentRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, middleware.ErrBadRequest(err.Error()))
		return
	}

	// Validate required fields
	if req.Content == "" {
		render.Render(w, r, middleware.ErrBadRequest("comment content is required"))
		return
	}

	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		render.Render(w, r, middleware.ErrUnauthorized("user not found in context"))
		return
	}

	// Create comment request
	commentReq := &dto.CommentRequest{
		Content: req.Content,
	}

	// Convert user role to entities.UserRole
	userRole := entities.UserRole(user.Role)

	// Execute use case
	result, err := h.addCommentUseCase.Execute(r.Context(), ticketID, user.ID, userRole, commentReq)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			render.Render(w, r, middleware.ErrNotFound("ticket not found"))
		case "user not authorized to comment":
			render.Render(w, r, middleware.ErrForbidden("user not authorized to comment on this ticket"))
		default:
			render.Render(w, r, middleware.ErrInternalServerError(err.Error()))
		}
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, result)
}

// GetComments handles retrieving all comments for a ticket
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketID")
	if ticketID == "" {
		render.Render(w, r, middleware.ErrBadRequest("ticket ID is required"))
		return
	}

	// Validate UUID
	if _, err := uuid.Parse(ticketID); err != nil {
		render.Render(w, r, middleware.ErrBadRequest("invalid ticket ID format"))
		return
	}

	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		render.Render(w, r, middleware.ErrUnauthorized("user not found in context"))
		return
	}

	// Create get comments request
	getReq := &dto.GetCommentsRequest{
		TicketID: ticketID,
		UserID:   user.ID,
		UserRole: user.Role,
	}

	// Execute use case
	result, err := h.getCommentsUseCase.Execute(r.Context(), getReq)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			render.Render(w, r, middleware.ErrNotFound("ticket not found"))
		case "user not authorized to view comments":
			render.Render(w, r, middleware.ErrForbidden("user not authorized to view comments on this ticket"))
		default:
			render.Render(w, r, middleware.ErrInternalServerError(err.Error()))
		}
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, result)
}