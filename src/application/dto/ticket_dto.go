package dto

import (
	"fmt"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// CreateTicketRequest represents a request to create a ticket
type CreateTicketRequest struct {
	Title          string `json:"title" validate:"required,max=255"`
	Description    string `json:"description" validate:"required"`
	Category       string `json:"category" validate:"required,oneof=office_supplies facility_maintenance pantry_supplies meeting_room office_furniture general_service"`
	Priority       string `json:"priority" validate:"required,oneof=low medium high"`
	EstimatedCost  int64  `json:"estimated_cost" validate:"required,min=0"`
	RequesterID    string `json:"requester_id" validate:"required"`
}

// Validate validates the CreateTicketRequest
func (req *CreateTicketRequest) Validate() error {
	if req.Title == "" {
		return fmt.Errorf("title is required")
	}
	if len(req.Title) > 255 {
		return fmt.Errorf("title must be 255 characters or less")
	}
	if req.Description == "" {
		return fmt.Errorf("description is required")
	}
	if req.EstimatedCost < 0 {
		return fmt.Errorf("estimated cost cannot be negative")
	}
	if req.RequesterID == "" {
		return fmt.Errorf("requester ID is required")
	}
	return nil
}

// GetTicketsRequest represents a request to get tickets with pagination
type GetTicketsRequest struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role" validate:"required,oneof=requester approver admin"`
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Status   string `json:"status,omitempty"`
	Category string `json:"category,omitempty"`
}

// Validate validates the GetTicketsRequest
func (req *GetTicketsRequest) Validate() error {
	if req.UserRole == "" {
		return fmt.Errorf("user role is required")
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	return nil
}

// GetTicketRequest represents a request to get a single ticket
type GetTicketRequest struct {
	TicketID string `json:"ticket_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	UserRole string `json:"user_role" validate:"required,oneof=requester approver admin"`
}

// Validate validates the GetTicketRequest
func (req *GetTicketRequest) Validate() error {
	if req.TicketID == "" {
		return fmt.Errorf("ticket ID is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if req.UserRole == "" {
		return fmt.Errorf("user role is required")
	}
	return nil
}

// UpdateTicketRequest represents a request to update a ticket
type UpdateTicketRequest struct {
	Title        *string `json:"title,omitempty" validate:"omitempty,max=255"`
	Description  *string `json:"description,omitempty"`
	Priority     *string `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	ActualCost   *int64  `json:"actual_cost,omitempty" validate:"omitempty,min=0"`
	Status       *string `json:"status,omitempty" validate:"omitempty,oneof=pending waiting_approval approved rejected in_progress completed closed"`
	Reason       *string `json:"reason,omitempty"`
	UpdatedBy    *string `json:"updated_by,omitempty"`
}

// Validate validates the UpdateTicketRequest
func (req *UpdateTicketRequest) Validate() error {
	if req.Title != nil && *req.Title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	if req.Title != nil && len(*req.Title) > 255 {
		return fmt.Errorf("title must be 255 characters or less")
	}
	if req.ActualCost != nil && *req.ActualCost < 0 {
		return fmt.Errorf("actual cost cannot be negative")
	}
	return nil
}

// AssignTicketRequest represents a request to assign a ticket
type AssignTicketRequest struct {
	AdminID string `json:"admin_id" validate:"required"`
}

// Validate validates the AssignTicketRequest
func (req *AssignTicketRequest) Validate() error {
	if req.AdminID == "" {
		return fmt.Errorf("admin ID is required")
	}
	return nil
}

// CommentRequest represents a request to add a comment
type CommentRequest struct {
	Content string `json:"content" validate:"required,max=1000"`
}

// Validate validates the CommentRequest
func (req *CommentRequest) Validate() error {
	if req.Content == "" {
		return fmt.Errorf("comment content is required")
	}
	if len(req.Content) > 1000 {
		return fmt.Errorf("comment must be 1000 characters or less")
	}
	return nil
}

// TicketResponse represents a ticket response
type TicketResponse struct {
	ID               string                 `json:"id"`
	TicketNumber     string                 `json:"ticket_number"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Category         string                 `json:"category"`
	Priority         string                 `json:"priority"`
	Status           string                 `json:"status"`
	RequesterID      string                 `json:"requester_id"`
	AssignedAdminID  *string                `json:"assigned_admin_id"`
	EstimatedCost    int64                  `json:"estimated_cost"`
	ActualCost       *int64                 `json:"actual_cost"`
	RequiresApproval bool                   `json:"requires_approval"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
	CompletedAt      *string                `json:"completed_at"`
	AssignedAt       *string                `json:"assigned_at"`
	Comments         []*CommentResponse     `json:"comments,omitempty"`
	StatusHistory    []*StatusHistoryResponse `json:"status_history,omitempty"`
}

// GetTicketsResponse represents a response for getting multiple tickets
type GetTicketsResponse struct {
	Tickets []*TicketResponse `json:"tickets"`
	Page    int               `json:"page"`
	Limit   int               `json:"limit"`
	Total   int               `json:"total"`
}

// CommentResponse represents a comment response
type CommentResponse struct {
	ID        string `json:"id"`
	TicketID  string `json:"ticket_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// StatusHistoryResponse represents a status history response
type StatusHistoryResponse struct {
	ID         string `json:"id"`
	TicketID   string `json:"ticket_id"`
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
	ChangedBy  string `json:"changed_by"`
	Comments   string `json:"comments"`
	CreatedAt  string `json:"created_at"`
}

// TicketFromEntity converts a domain entity to DTO
func TicketFromEntity(ticket *entities.Ticket) *TicketResponse {
	response := &TicketResponse{
		ID:               ticket.GetID(),
		TicketNumber:     ticket.GetTicketNumber(),
		Title:            ticket.GetTitle(),
		Description:      ticket.GetDescription(),
		Category:         string(ticket.GetCategory()),
		Priority:         string(ticket.GetPriority()),
		Status:           string(ticket.GetStatus()),
		RequesterID:      ticket.GetRequesterID(),
		EstimatedCost:    ticket.GetEstimatedCost().Amount,
		RequiresApproval: ticket.RequiresApproval(),
		CreatedAt:        ticket.GetCreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:        ticket.GetUpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	// Handle optional fields
	if ticket.GetActualCost() != nil {
		response.ActualCost = &ticket.GetActualCost().Amount
	}

	if ticket.GetAssignedAdminID() != nil {
		response.AssignedAdminID = ticket.GetAssignedAdminID()
	}

	if ticket.GetCompletedAt() != nil {
		completedAt := ticket.GetCompletedAt().Format("2006-01-02T15:04:05Z")
		response.CompletedAt = &completedAt
	}

	if ticket.GetAssignedAt() != nil {
		assignedAt := ticket.GetAssignedAt().Format("2006-01-02T15:04:05Z")
		response.AssignedAt = &assignedAt
	}

	// Convert comments
	comments := ticket.GetComments()
	response.Comments = make([]*CommentResponse, len(comments))
	for i, comment := range comments {
		response.Comments[i] = CommentFromEntity(comment)
	}

	// Convert status history
	history := ticket.GetStatusHistory()
	response.StatusHistory = make([]*StatusHistoryResponse, len(history))
	for i, h := range history {
		response.StatusHistory[i] = StatusHistoryFromEntity(h)
	}

	return response
}

// CommentFromEntity converts a domain comment to DTO
func CommentFromEntity(comment *entities.Comment) *CommentResponse {
	return &CommentResponse{
		ID:        comment.GetID(),
		TicketID:  comment.GetTicketID(),
		UserID:    comment.GetUserID(),
		Content:   comment.GetContent(),
		CreatedAt: comment.GetCreatedAt().Format("2006-01-02T15:04:05Z"),
	}
}

// StatusHistoryFromEntity converts domain status history to DTO
func StatusHistoryFromEntity(history *entities.StatusHistory) *StatusHistoryResponse {
	return &StatusHistoryResponse{
		ID:         history.ID,
		TicketID:   history.TicketID,
		FromStatus: string(history.FromStatus),
		ToStatus:   string(history.ToStatus),
		ChangedBy:  history.ChangedBy,
		Comments:   history.Comments,
		CreatedAt:  history.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// MoneyFromValueObject converts a Money value object to a simple response
func MoneyFromValueObject(money *valueobjects.Money) map[string]interface{} {
	return map[string]interface{}{
		"amount":   money.Amount,
		"currency": money.Currency,
		"display":  money.FormatIndonesian(),
	}
}