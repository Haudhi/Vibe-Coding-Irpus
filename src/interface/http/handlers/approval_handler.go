package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/application/usecases"
	"github.com/company/ga-ticketing/src/interface/http/middleware"
)

// ApprovalHandler handles HTTP requests for approval operations
type ApprovalHandler struct {
	approveUseCase *usecases.ApproveTicketUseCase
	rejectUseCase  *usecases.RejectTicketUseCase
}

// NewApprovalHandler creates a new approval handler
func NewApprovalHandler(
	approveUseCase *usecases.ApproveTicketUseCase,
	rejectUseCase *usecases.RejectTicketUseCase,
) *ApprovalHandler {
	return &ApprovalHandler{
		approveUseCase: approveUseCase,
		rejectUseCase:  rejectUseCase,
	}
}

// ApproveRequest represents the request body for approving a ticket
type ApproveRequest struct {
	Comments string `json:"comments" validate:"max=1000"`
}

// RejectRequest represents the request body for rejecting a ticket
type RejectRequest struct {
	Comments string `json:"comments" validate:"required,max=1000"`
}

// ApproveTicket handles the approval of a ticket
func (h *ApprovalHandler) ApproveTicket(w http.ResponseWriter, r *http.Request) {
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

	var req ApproveRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, middleware.ErrBadRequest(err.Error()))
		return
	}

	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		render.Render(w, r, middleware.ErrUnauthorized("user not found in context"))
		return
	}

	// Create approval request
	approvalReq := &dto.ApproveTicketRequest{
		TicketID: ticketID,
		ApproverID: user.ID,
		Comments: req.Comments,
	}

	// Execute use case
	result, err := h.approveUseCase.Execute(r.Context(), approvalReq)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			render.Render(w, r, middleware.ErrNotFound("ticket not found"))
		case "ticket does not require approval":
			render.Render(w, r, middleware.ErrBadRequest("ticket does not require approval"))
		case "already approved":
			render.Render(w, r, middleware.ErrConflict("ticket has already been approved"))
		case "already rejected":
			render.Render(w, r, middleware.ErrConflict("ticket has already been rejected"))
		case "user not authorized to approve":
			render.Render(w, r, middleware.ErrForbidden("user not authorized to approve this ticket"))
		default:
			render.Render(w, r, middleware.ErrInternalServerError(err.Error()))
		}
		return
	}

	render.JSON(w, http.StatusOK, result)
}

// RejectTicket handles the rejection of a ticket
func (h *ApprovalHandler) RejectTicket(w http.ResponseWriter, r *http.Request) {
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

	var req RejectRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Render(w, r, middleware.ErrBadRequest(err.Error()))
		return
	}

	// Validate required fields
	if req.Comments == "" {
		render.Render(w, r, middleware.ErrBadRequest("rejection reason is required"))
		return
	}

	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		render.Render(w, r, middleware.ErrUnauthorized("user not found in context"))
		return
	}

	// Create rejection request
	rejectionReq := &dto.RejectTicketRequest{
		TicketID: ticketID,
		ApproverID: user.ID,
		Comments: req.Comments,
	}

	// Execute use case
	result, err := h.rejectUseCase.Execute(r.Context(), rejectionReq)
	if err != nil {
		switch err.Error() {
		case "ticket not found":
			render.Render(w, r, middleware.ErrNotFound("ticket not found"))
		case "ticket does not require approval":
			render.Render(w, r, middleware.ErrBadRequest("ticket does not require approval"))
		case "already approved":
			render.Render(w, r, middleware.ErrConflict("ticket has already been approved"))
		case "already rejected":
			render.Render(w, r, middleware.ErrConflict("ticket has already been rejected"))
		case "user not authorized to reject":
			render.Render(w, r, middleware.ErrForbidden("user not authorized to reject this ticket"))
		default:
			render.Render(w, r, middleware.ErrInternalServerError(err.Error()))
		}
		return
	}

	render.JSON(w, http.StatusOK, result)
}

// GetPendingApprovals handles retrieving all pending approvals for the current user
func (h *ApprovalHandler) GetPendingApprovals(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		render.Render(w, r, middleware.ErrUnauthorized("user not found in context"))
		return
	}

	// Check if user is approver or admin
	if user.Role != "approver" && user.Role != "admin" {
		render.Render(w, r, middleware.ErrForbidden("user not authorized to view pending approvals"))
		return
	}

	// This would typically use a GetPendingApprovals use case
	// For now, return an empty list
	render.JSON(w, http.StatusOK, map[string]interface{}{
		"pending_approvals": []interface{}{},
		"count": 0,
	})
}