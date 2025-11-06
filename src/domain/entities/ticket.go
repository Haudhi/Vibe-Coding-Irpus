package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// TicketCategory represents the category of a ticket
type TicketCategory string

const (
	CategoryOfficeSupplies     TicketCategory = "office_supplies"
	CategoryFacilityMaintenance TicketCategory = "facility_maintenance"
	CategoryPantrySupplies     TicketCategory = "pantry_supplies"
	CategoryMeetingRoom        TicketCategory = "meeting_room"
	CategoryOfficeFurniture    TicketCategory = "office_furniture"
	CategoryGeneralService     TicketCategory = "general_service"
)

// TicketPriority represents the priority level of a ticket
type TicketPriority string

const (
	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
)

// TicketStatus represents the status of a ticket
type TicketStatus string

const (
	StatusPending          TicketStatus = "pending"
	StatusWaitingApproval  TicketStatus = "waiting_approval"
	StatusApproved         TicketStatus = "approved"
	StatusRejected         TicketStatus = "rejected"
	StatusInProgress       TicketStatus = "in_progress"
	StatusCompleted        TicketStatus = "completed"
	StatusClosed           TicketStatus = "closed"
)

// StatusHistory represents a change in ticket status
type StatusHistory struct {
	ID          string        `json:"id"`
	TicketID    string        `json:"ticket_id"`
	FromStatus  TicketStatus  `json:"from_status"`
	ToStatus    TicketStatus  `json:"to_status"`
	ChangedBy   string        `json:"changed_by"`
	Comments    string        `json:"comments"`
	CreatedAt   time.Time     `json:"created_at"`
}

// NewStatusHistory creates a new status history entry
func NewStatusHistory(ticketID string, fromStatus, toStatus TicketStatus, changedBy, comments string) *StatusHistory {
	return &StatusHistory{
		ID:        uuid.New().String(),
		TicketID:  ticketID,
		FromStatus: fromStatus,
		ToStatus:  toStatus,
		ChangedBy: changedBy,
		Comments:  comments,
		CreatedAt: time.Now(),
	}
}

// Ticket represents a service request ticket
type Ticket struct {
	id              string
	ticketNumber    string
	title           string
	description     string
	category        TicketCategory
	priority        TicketPriority
	status          TicketStatus
	requesterID     string
	assignedAdminID *string
	estimatedCost   *valueobjects.Money
	actualCost      *valueobjects.Money
	requiresApproval bool
	createdAt       time.Time
	updatedAt       time.Time
	completedAt     *time.Time
	assignedAt      *time.Time
	statusHistory   []*StatusHistory
	comments        []*Comment
}

// NewTicket creates a new ticket
func NewTicket(
	title, description string,
	category TicketCategory,
	priority TicketPriority,
	estimatedCost *valueobjects.Money,
	requesterID string,
) (*Ticket, error) {
	// Validate input
	if title == "" {
		return nil, errors.New("title is required")
	}
	if len(title) > 255 {
		return nil, errors.New("title must be 255 characters or less")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}
	if estimatedCost == nil {
		return nil, errors.New("estimated cost is required")
	}
	if estimatedCost.Amount < 0 {
		return nil, errors.New("estimated cost cannot be negative")
	}
	if requesterID == "" {
		return nil, errors.New("requester ID is required")
	}

	// Check if approval is required
	requiresApproval := category == CategoryOfficeFurniture || estimatedCost.Amount >= 500000

	now := time.Now()
	ticket := &Ticket{
		id:               uuid.New().String(),
		ticketNumber:     generateTicketNumber(now),
		title:            title,
		description:      description,
		category:         category,
		priority:         priority,
		status:           StatusPending,
		requesterID:      requesterID,
		estimatedCost:    estimatedCost,
		requiresApproval: requiresApproval,
		createdAt:        now,
		updatedAt:        now,
		statusHistory:    make([]*StatusHistory, 0),
		comments:         make([]*Comment, 0),
	}

	// Add initial status history
	ticket.statusHistory = append(ticket.statusHistory,
		NewStatusHistory(ticket.id, "", StatusPending, requesterID, "Ticket created"),
	)

	// Set initial status to waiting_approval if approval is required
	if requiresApproval {
		ticket.status = StatusWaitingApproval
		ticket.statusHistory[0].ToStatus = StatusWaitingApproval
	}

	return ticket, nil
}

// Getters
func (t *Ticket) GetID() string                      { return t.id }
func (t *Ticket) GetTicketNumber() string           { return t.ticketNumber }
func (t *Ticket) GetTitle() string                  { return t.title }
func (t *Ticket) GetDescription() string            { return t.description }
func (t *Ticket) GetCategory() TicketCategory       { return t.category }
func (t *Ticket) GetPriority() TicketPriority       { return t.priority }
func (t *Ticket) GetStatus() TicketStatus           { return t.status }
func (t *Ticket) GetRequesterID() string            { return t.requesterID }
func (t *Ticket) GetAssignedAdminID() *string       { return t.assignedAdminID }
func (t *Ticket) GetEstimatedCost() *valueobjects.Money  { return t.estimatedCost }
func (t *Ticket) GetActualCost() *valueobjects.Money     { return t.actualCost }
func (t *Ticket) GetCreatedAt() time.Time           { return t.createdAt }
func (t *Ticket) GetUpdatedAt() time.Time           { return t.updatedAt }
func (t *Ticket) GetCompletedAt() *time.Time        { return t.completedAt }
func (t *Ticket) GetAssignedAt() *time.Time         { return t.assignedAt }
func (t *Ticket) GetStatusHistory() []*StatusHistory { return t.statusHistory }
func (t *Ticket) GetComments() []*Comment           { return t.comments }
func (t *Ticket) RequiresApproval() bool            { return t.requiresApproval }

// Setters (business logic)
func (t *Ticket) SetStatus(status TicketStatus, reason, updatedBy string) error {
	if updatedBy == "" {
		return errors.New("updated by is required")
	}

	if !t.canTransitionTo(status) {
		return fmt.Errorf("invalid status transition from %s to %s", t.status, status)
	}

	fromStatus := t.status
	t.status = status
	t.updatedAt = time.Now()

	// Add to status history
	t.statusHistory = append(t.statusHistory,
		NewStatusHistory(t.id, fromStatus, status, updatedBy, reason),
	)

	// Set completed timestamp if status is completed
	if status == StatusCompleted {
		now := time.Now()
		t.completedAt = &now
	}

	return nil
}

func (t *Ticket) SetTitle(title string) error {
	if title == "" {
		return errors.New("title is required")
	}
	if len(title) > 255 {
		return errors.New("title must be 255 characters or less")
	}

	t.title = title
	t.updatedAt = time.Now()
	return nil
}

func (t *Ticket) SetDescription(description string) error {
	if description == "" {
		return errors.New("description is required")
	}

	t.description = description
	t.updatedAt = time.Now()
	return nil
}

func (t *Ticket) SetPriority(priority TicketPriority) {
	t.priority = priority
	t.updatedAt = time.Now()
}

func (t *Ticket) SetEstimatedCost(cost *valueobjects.Money) error {
	if cost == nil {
		return errors.New("estimated cost is required")
	}
	if cost.Amount < 0 {
		return errors.New("estimated cost cannot be negative")
	}

	t.estimatedCost = cost
	t.updatedAt = time.Now()

	// Re-check approval requirement
	t.requiresApproval = t.category == CategoryOfficeFurniture || cost.Amount >= 500000

	return nil
}

func (t *Ticket) SetActualCost(cost *valueobjects.Money) error {
	if cost == nil {
		return errors.New("actual cost is required")
	}
	if cost.Amount < 0 {
		return errors.New("actual cost cannot be negative")
	}

	t.actualCost = cost
	t.updatedAt = time.Now()
	return nil
}

func (t *Ticket) AssignToAdmin(adminID string) error {
	if adminID == "" {
		return errors.New("admin ID is required")
	}

	if t.assignedAdminID != nil && *t.assignedAdminID != "" {
		return errors.New("ticket is already assigned")
	}

	t.assignedAdminID = &adminID
	now := time.Now()
	t.assignedAt = &now

	// Update status to in_progress
	if t.status == StatusPending || t.status == StatusApproved {
		t.SetStatus(StatusInProgress, "Assigned to admin", adminID)
	}

	return nil
}

func (t *Ticket) ReassignToAdmin(adminID string) error {
	if adminID == "" {
		return errors.New("admin ID is required")
	}

	t.assignedAdminID = &adminID
	t.updatedAt = time.Now()
	return nil
}

func (t *Ticket) Unassign() {
	t.assignedAdminID = nil
	t.updatedAt = time.Now()
}

func (t *Ticket) AddComment(content, userID string) (*Comment, error) {
	if content == "" {
		return nil, errors.New("comment content is required")
	}
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	comment := NewComment(t.id, content, userID)
	t.comments = append(t.comments, comment)
	t.updatedAt = time.Now()

	return comment, nil
}

func (t *Ticket) GetTimeInCurrentStatus() time.Duration {
	if len(t.statusHistory) == 0 {
		return 0
	}

	// Find the last status change
	for i := len(t.statusHistory) - 1; i >= 0; i-- {
		if t.statusHistory[i].ToStatus == t.status {
			return time.Since(t.statusHistory[i].CreatedAt)
		}
	}

	// If no status history found, return time since creation
	return time.Since(t.createdAt)
}

func (t *Ticket) CanBeViewedBy(userID, userRole string) bool {
	// Admin can view any ticket
	if userRole == string(RoleAdmin) {
		return true
	}

	// User can view their own tickets
	if t.requesterID == userID {
		return true
	}

	// Assigned admin can view the ticket
	if t.assignedAdminID != nil && *t.assignedAdminID == userID {
		return true
	}

	// Approver can view tickets requiring approval
	if userRole == string(RoleApprover) && t.requiresApproval {
		return true
	}

	return false
}

func (t *Ticket) canTransitionTo(newStatus TicketStatus) bool {
	if t.status == newStatus {
		return true // No change needed
	}

	// Define valid transitions
	validTransitions := map[TicketStatus][]TicketStatus{
		StatusPending:         {StatusWaitingApproval, StatusInProgress, StatusClosed},
		StatusWaitingApproval: {StatusApproved, StatusRejected, StatusClosed},
		StatusApproved:        {StatusInProgress, StatusClosed},
		StatusRejected:        {StatusClosed},
		StatusInProgress:      {StatusCompleted, StatusClosed},
		StatusCompleted:       {StatusClosed},
		StatusClosed:          {}, // No transitions from closed
	}

	allowed, exists := validTransitions[t.status]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}

	return false
}

// generateTicketNumber generates a unique ticket number
func generateTicketNumber(createdAt time.Time) string {
	year := createdAt.Year()
	sequence := 1 // In a real implementation, this would come from a sequence generator
	return fmt.Sprintf("GA-%d-%04d", year, sequence)
}

// Comment represents a comment on a ticket
type Comment struct {
	id        string
	ticketID  string
	userID    string
	content   string
	createdAt time.Time
}

// NewComment creates a new comment
func NewComment(ticketID, content, userID string) *Comment {
	return &Comment{
		id:        uuid.New().String(),
		ticketID:  ticketID,
		userID:    userID,
		content:   content,
		createdAt: time.Now(),
	}
}

// Getters for Comment
func (c *Comment) GetID() string        { return c.id }
func (c *Comment) GetTicketID() string  { return c.ticketID }
func (c *Comment) GetUserID() string    { return c.userID }
func (c *Comment) GetContent() string   { return c.content }
func (c *Comment) GetCreatedAt() time.Time { return c.createdAt }