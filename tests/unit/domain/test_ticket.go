package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

func TestTicket_NewTicket(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		category    entities.TicketCategory
		priority    entities.TicketPriority
		estimatedCost int64
		requesterID string
		wantErr     bool
		errMessage  string
	}{
		{
			name:          "valid ticket",
			title:         "Request for office supplies",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: 250000,
			requesterID:   uuid.New().String(),
			wantErr:       false,
		},
		{
			name:          "empty title",
			title:         "",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: 250000,
			requesterID:   uuid.New().String(),
			wantErr:       true,
			errMessage:    "title is required",
		},
		{
			name:          "title too long",
			title:         string(make([]byte, 256)), // 256 characters
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: 250000,
			requesterID:   uuid.New().String(),
			wantErr:       true,
			errMessage:    "title must be 255 characters or less",
		},
		{
			name:          "negative estimated cost",
			title:         "Request for office supplies",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: -1000,
			requesterID:   uuid.New().String(),
			wantErr:       true,
			errMessage:    "estimated cost cannot be negative",
		},
		{
			name:          "office furniture requires approval",
			title:         "New office chair",
			description:   "Ergonomic chair for back pain",
			category:      entities.CategoryOfficeFurniture,
			priority:      entities.PriorityHigh,
			estimatedCost: 300000,
			requesterID:   uuid.New().String(),
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket, err := entities.NewTicket(
				tt.title,
				tt.description,
				tt.category,
				tt.priority,
				valueobjects.NewMoney(tt.estimatedCost),
				tt.requesterID,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
				assert.Nil(t, ticket)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, ticket)
				assert.Equal(t, tt.title, ticket.GetTitle())
				assert.Equal(t, tt.description, ticket.GetDescription())
				assert.Equal(t, tt.category, ticket.GetCategory())
				assert.Equal(t, tt.priority, ticket.GetPriority())
				assert.Equal(t, tt.estimatedCost, ticket.GetEstimatedCost().Amount)
				assert.Equal(t, tt.requesterID, ticket.GetRequesterID())
				assert.Equal(t, entities.StatusPending, ticket.GetStatus())
				assert.False(t, ticket.RequiresApproval())

				// Office furniture should always require approval
				if tt.category == entities.CategoryOfficeFurniture {
					assert.True(t, ticket.RequiresApproval())
				}

				// High cost (>=500,000) should require approval
				if tt.estimatedCost >= 500000 {
					assert.True(t, ticket.RequiresApproval())
				}
			}
		})
	}
}

func TestTicket_AssignToAdmin(t *testing.T) {
	ticket := createTestTicket(t)
	adminID := uuid.New().String()

	err := ticket.AssignToAdmin(adminID)
	assert.NoError(t, err)
	assert.Equal(t, adminID, ticket.GetAssignedAdminID())
	assert.Equal(t, entities.StatusInProgress, ticket.GetStatus())
	assert.NotNil(t, ticket.GetAssignedAt())

	// Test assignment with empty admin ID
	err = ticket.AssignToAdmin("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "admin ID is required")
}

func TestTicket_UpdateStatus(t *testing.T) {
	ticket := createTestTicket(t)

	tests := []struct {
		name        string
		fromStatus  entities.TicketStatus
		toStatus    entities.TicketStatus
		reason      string
		updatedBy   string
		wantErr     bool
		errMessage  string
	}{
		{
			name:       "valid transition: pending to waiting_approval",
			fromStatus: entities.StatusPending,
			toStatus:   entities.StatusWaitingApproval,
			reason:     "High cost item requires approval",
			updatedBy:  uuid.New().String(),
			wantErr:    false,
		},
		{
			name:       "valid transition: pending to in_progress",
			fromStatus: entities.StatusPending,
			toStatus:   entities.StatusInProgress,
			reason:     "Direct processing for low cost items",
			updatedBy:  uuid.New().String(),
			wantErr:    false,
		},
		{
			name:       "valid transition: waiting_approval to approved",
			fromStatus: entities.StatusWaitingApproval,
			toStatus:   entities.StatusApproved,
			reason:     "Budget approved",
			updatedBy:  uuid.New().String(),
			wantErr:    false,
		},
		{
			name:       "valid transition: approved to in_progress",
			fromStatus: entities.StatusApproved,
			toStatus:   entities.StatusInProgress,
			reason:     "Starting work on approved ticket",
			updatedBy:  uuid.New().String(),
			wantErr:    false,
		},
		{
			name:       "valid transition: in_progress to completed",
			fromStatus: entities.StatusInProgress,
			toStatus:   entities.StatusCompleted,
			reason:     "Work completed successfully",
			updatedBy:  uuid.New().String(),
			wantErr:    false,
		},
		{
			name:       "valid transition: completed to closed",
			fromStatus: entities.StatusCompleted,
			toStatus:   entities.StatusClosed,
			reason:     "Ticket closed",
			updatedBy:  uuid.New().String(),
			wantErr:    false,
		},
		{
			name:       "invalid transition: pending to completed",
			fromStatus: entities.StatusPending,
			toStatus:   entities.StatusCompleted,
			reason:     "Invalid transition",
			updatedBy:  uuid.New().String(),
			wantErr:    true,
			errMessage: "invalid status transition",
		},
		{
			name:       "empty updated by",
			fromStatus: entities.StatusPending,
			toStatus:   entities.StatusInProgress,
			reason:     "No updater specified",
			updatedBy:  "",
			wantErr:    true,
			errMessage: "updated by is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset ticket status
			ticket.SetStatus(tt.fromStatus, "", "")

			// TODO: UpdateStatus method needs to be implemented on Ticket entity
			// err := ticket.UpdateStatus(tt.toStatus, tt.reason, tt.updatedBy)
			var err error

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.toStatus, ticket.GetStatus())

				// Check if status history was added
				if tt.toStatus != tt.fromStatus {
					history := ticket.GetStatusHistory()
					assert.Len(t, history, 1)
					lastHistory := history[len(history)-1]
					assert.Equal(t, tt.fromStatus, lastHistory.FromStatus)
					assert.Equal(t, tt.toStatus, lastHistory.ToStatus)
					assert.Equal(t, tt.reason, lastHistory.Comments)
					assert.Equal(t, tt.updatedBy, lastHistory.ChangedBy)
				}

				// Check completion timestamp
				if tt.toStatus == entities.StatusCompleted {
					assert.NotNil(t, ticket.GetCompletedAt())
				}
			}
		})
	}
}

func TestTicket_SetActualCost(t *testing.T) {
	ticket := createTestTicket(t)
	actualCost := int64(300000)

	err := ticket.SetActualCost(valueobjects.NewMoney(actualCost))
	assert.NoError(t, err)
	assert.Equal(t, actualCost, ticket.GetActualCost().Amount)

	// Test negative actual cost
	err = ticket.SetActualCost(valueobjects.NewMoney(-1000))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "actual cost cannot be negative")
}

func TestTicket_AddComment(t *testing.T) {
	ticket := createTestTicket(t)
	userID := uuid.New().String()
	content := "Please update the status"

	comment, err := ticket.AddComment(content, userID)
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, content, comment.GetContent())
	assert.Equal(t, userID, comment.GetUserID())
	assert.Equal(t, ticket.GetID(), comment.GetTicketID())

	// Check comment was added to ticket
	comments := ticket.GetComments()
	assert.Len(t, comments, 1)
	assert.Equal(t, comment, comments[0])

	// Test empty comment
	_, err = ticket.AddComment("", userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "comment content is required")
}

func TestTicket_CanBeViewedBy(t *testing.T) {
	ticket := createTestTicket(t)
	requesterID := ticket.GetRequesterID()
	adminID := uuid.New().String()
	otherUserID := uuid.New().String()

	// Ticket should be viewable by requester
	assert.True(t, ticket.CanBeViewedBy(requesterID, "requester"))
	assert.True(t, ticket.CanBeViewedBy(requesterID, "admin"))
	assert.True(t, ticket.CanBeViewedBy(requesterID, "approver"))

	// Admin can view any ticket
	assert.True(t, ticket.CanBeViewedBy(adminID, "admin"))

	// Approver can view tickets requiring approval
	// TODO: UpdateStatus method needs to be implemented on Ticket entity
	// ticket.UpdateStatus(entities.StatusWaitingApproval, "High cost", uuid.New().String())
	ticket.SetStatus(entities.StatusWaitingApproval, "High cost", uuid.New().String())
	assert.True(t, ticket.CanBeViewedBy(otherUserID, "approver"))

	// Other requester cannot view this ticket
	assert.False(t, ticket.CanBeViewedBy(otherUserID, "requester"))
}

func TestTicket_RequiresApproval(t *testing.T) {
	tests := []struct {
		name          string
		category      entities.TicketCategory
		estimatedCost int64
		expected      bool
	}{
		{
			name:          "office supplies low cost",
			category:      entities.CategoryOfficeSupplies,
			estimatedCost: 100000,
			expected:      false,
		},
		{
			name:          "office supplies high cost",
			category:      entities.CategoryOfficeSupplies,
			estimatedCost: 500000,
			expected:      true,
		},
		{
			name:          "office furniture low cost",
			category:      entities.CategoryOfficeFurniture,
			estimatedCost: 100000,
			expected:      true,
		},
		{
			name:          "meeting room setup",
			category:      entities.CategoryMeetingRoom,
			estimatedCost: 300000,
			expected:      false,
		},
		{
			name:          "facility maintenance high cost",
			category:      entities.CategoryFacilityMaintenance,
			estimatedCost: 600000,
			expected:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket, err := entities.NewTicket(
				"Test ticket",
				"Test description",
				tt.category,
			 entities.PriorityMedium,
				valueobjects.NewMoney(tt.estimatedCost),
				uuid.New().String(),
			)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, ticket.RequiresApproval())
		})
	}
}

func TestTicket_GetTimeInStatus(t *testing.T) {
	ticket := createTestTicket(t)

	// Initially in pending status
	timeInPending := ticket.GetTimeInCurrentStatus()
	assert.True(t, timeInPending >= 0)

	// Wait a bit and check again
	time.Sleep(10 * time.Millisecond)
	timeInPending2 := ticket.GetTimeInCurrentStatus()
	assert.True(t, timeInPending2 > timeInPending)
}

func TestTicket_GenerateTicketNumber(t *testing.T) {
	// Test that ticket numbers are generated correctly
	ticket1 := createTestTicket(t)
	ticket2 := createTestTicket(t)

	num1 := ticket1.GetTicketNumber()
	num2 := ticket2.GetTicketNumber()

	assert.NotEmpty(t, num1)
	assert.NotEmpty(t, num2)
	assert.NotEqual(t, num1, num2)

	// Ticket numbers should follow the format GA-YYYY-NNNN
	assert.Regexp(t, `^GA-\d{4}-\d{4}$`, num1)
	assert.Regexp(t, `^GA-\d{4}-\d{4}$`, num2)
}

// Helper function to create a test ticket
func createTestTicket(t *testing.T) *entities.Ticket {
	ticket, err := entities.NewTicket(
		"Test ticket",
		"This is a test ticket",
		entities.CategoryOfficeSupplies,
		entities.PriorityMedium,
		valueobjects.NewMoney(250000),
		uuid.New().String(),
	)
	require.NoError(t, err)
	return ticket
}