package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/application/usecases"
	"github.com/company/ga-ticketing/src/domain/entities"
)

// TicketHandler handles ticket-related HTTP requests
type TicketHandler struct {
	createTicketUC  *usecases.CreateTicketUseCase
	getTicketsUC    *usecases.GetTicketsUseCase
	getTicketUC     *usecases.GetTicketUseCase
	updateTicketUC  *usecases.UpdateTicketUseCase
	assignTicketUC  *usecases.AssignTicketUseCase
	getCommentsUC   *usecases.GetCommentsUseCase
	addCommentUC    *usecases.AddCommentUseCase
	approveTicketUC *usecases.ApproveTicketUseCase
	rejectTicketUC  *usecases.RejectTicketUseCase
	logger          *zap.Logger
}

// NewTicketHandler creates a new TicketHandler
func NewTicketHandler(
	createTicketUC *usecases.CreateTicketUseCase,
	getTicketsUC *usecases.GetTicketsUseCase,
	getTicketUC *usecases.GetTicketUseCase,
	updateTicketUC *usecases.UpdateTicketUseCase,
	assignTicketUC *usecases.AssignTicketUseCase,
	getCommentsUC *usecases.GetCommentsUseCase,
	addCommentUC *usecases.AddCommentUseCase,
	approveTicketUC *usecases.ApproveTicketUseCase,
	rejectTicketUC *usecases.RejectTicketUseCase,
	logger *zap.Logger,
) *TicketHandler {
	return &TicketHandler{
		createTicketUC:  createTicketUC,
		getTicketsUC:    getTicketsUC,
		getTicketUC:     getTicketUC,
		updateTicketUC:  updateTicketUC,
		assignTicketUC:  assignTicketUC,
		getCommentsUC:   getCommentsUC,
		addCommentUC:    addCommentUC,
		approveTicketUC: approveTicketUC,
		rejectTicketUC:  rejectTicketUC,
		logger:          logger,
	}
}

// CreateTicket handles POST /v1/tickets
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTicketRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Warn("Failed to decode ticket request", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Get user info from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}
	_ = userRole

	req.RequesterID = userID

	ticket, err := h.createTicketUC.Execute(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create ticket",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Ticket created successfully",
		zap.String("ticket_id", ticket.ID),
		zap.String("ticket_number", ticket.TicketNumber),
		zap.String("user_id", userID),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, ticket)
}

// GetTickets handles GET /v1/tickets
func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	_ = page
	_ = limit

	// Get user info from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	req := &dto.GetTicketsRequest{
		UserID:   userID,
		UserRole: userRole,
		Page:     page,
		Limit:    limit,
		Status:   r.URL.Query().Get("status"),
		Category: r.URL.Query().Get("category"),
	}

	response, err := h.getTicketsUC.Execute(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get tickets",
			zap.String("user_id", userID),
			zap.String("user_role", userRole),
			zap.Error(err),
		)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// GetTicket handles GET /v1/tickets/{ticketId}
func (h *TicketHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	// Get user info from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	req := &dto.GetTicketRequest{
		TicketID: ticketID,
		UserID:   userID,
		UserRole: userRole,
	}

	ticket, err := h.getTicketUC.Execute(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get ticket",
			zap.String("ticket_id", ticketID),
			zap.String("user_id", userID),
			zap.Error(err),
		)

		if err.Error() == "ticket not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Ticket not found"})
			return
		}

		if err.Error() == "access denied" {
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, map[string]string{"error": "Access denied"})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, ticket)
}

// UpdateTicket handles PUT /v1/tickets/{ticketId}
func (h *TicketHandler) UpdateTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	var req dto.UpdateTicketRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Get user info from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRoleStr, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	userRole, err := entities.RoleFromString(userRoleStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user role"})
		return
	}

	ticket, err := h.updateTicketUC.Execute(r.Context(), ticketID, userID, userRole, &req)
	if err != nil {
		h.logger.Error("Failed to update ticket",
			zap.String("ticket_id", ticketID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Ticket updated successfully",
		zap.String("ticket_id", ticketID),
		zap.String("user_id", userID),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, ticket)
}

// AssignTicket handles POST /v1/tickets/{ticketId}/assign
func (h *TicketHandler) AssignTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	var req dto.AssignTicketRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Get user info from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRoleStr, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	userRole, err := entities.RoleFromString(userRoleStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user role"})
		return
	}

	ticket, err := h.assignTicketUC.Execute(r.Context(), ticketID, userID, userRole, &req)
	if err != nil {
		h.logger.Error("Failed to assign ticket",
			zap.String("ticket_id", ticketID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Ticket assigned successfully",
		zap.String("ticket_id", ticketID),
		zap.String("admin_id", req.AdminID),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, ticket)
}

// AddComment handles POST /v1/tickets/{ticketId}/comments
func (h *TicketHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	var req dto.CommentRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Get user info from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRoleStr, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	userRole, err := entities.RoleFromString(userRoleStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid user role"})
		return
	}

	comment, err := h.addCommentUC.Execute(r.Context(), ticketID, userID, userRole, &req)
	if err != nil {
		h.logger.Error("Failed to add comment",
			zap.String("ticket_id", ticketID),
			zap.String("user_id", userID),
			zap.Error(err),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Comment added successfully",
		zap.String("ticket_id", ticketID),
		zap.String("user_id", userID),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, comment)
}

// GetComments handles GET /v1/tickets/{ticketId}/comments
func (h *TicketHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "ticketId")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	// Get pagination parameters
	page := 1
	limit := 50

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get user info from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	req := &dto.GetCommentsRequest{
		TicketID: ticketID,
		UserID:   userID,
		UserRole: userRole,
		Page:     page,
		Limit:    limit,
	}

	response, err := h.getCommentsUC.Execute(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get comments",
			zap.String("ticket_id", ticketID),
			zap.String("user_id", userID),
			zap.Error(err),
		)

		if err.Error() == "access denied or ticket not found: ticket not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Ticket not found"})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}