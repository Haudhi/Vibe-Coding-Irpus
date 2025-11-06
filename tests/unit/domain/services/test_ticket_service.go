package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// MockTicketRepository is a mock implementation of TicketRepository
type MockTicketRepository struct {
	mock.Mock
}

func (m *MockTicketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) GetByID(ctx context.Context, id string) (*entities.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*entities.Ticket, error) {
	args := m.Called(ctx, ticketNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetByRequesterID(ctx context.Context, requesterID string, limit, offset int) ([]*entities.Ticket, error) {
	args := m.Called(ctx, requesterID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Ticket), args.Error(1)
}

func (m *MockTicketRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Ticket, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Ticket), args.Error(1)
}

func (m *MockTicketRepository) Update(ctx context.Context, ticket *entities.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTicketRepository) GetNextSequenceNumber(ctx context.Context, year int) (int, error) {
	args := m.Called(ctx, year)
	return args.Int(0), args.Error(1)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestTicketService_CreateTicket(t *testing.T) {
	ctx := context.Background()
	mockTicketRepo := &MockTicketRepository{}
	mockUserRepo := &MockUserRepository{}
	service := services.NewTicketService(mockTicketRepo, mockUserRepo)

	requesterID := uuid.New().String()
	validUser := createTestUser(t, requesterID)

	tests := []struct {
		name          string
		title         string
		description   string
		category      entities.TicketCategory
		priority      entities.TicketPriority
		estimatedCost int64
		requesterID   string
		setupMocks    func()
		wantErr       bool
		errMessage    string
	}{
		{
			name:          "valid ticket creation",
			title:         "Request for office supplies",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: 250000,
			requesterID:   requesterID,
			setupMocks: func() {
				mockUserRepo.On("GetByID", ctx, requesterID).Return(validUser, nil)
				mockTicketRepo.On("Create", ctx, mock.AnythingOfType("*entities.Ticket")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "high cost ticket requires approval",
			title:         "New office furniture",
			description:   "Executive desk and chair",
			category:      entities.CategoryOfficeFurniture,
			priority:      entities.PriorityHigh,
			estimatedCost: 1500000,
			requesterID:   requesterID,
			setupMocks: func() {
				mockUserRepo.On("GetByID", ctx, requesterID).Return(validUser, nil)
				mockTicketRepo.On("Create", ctx, mock.AnythingOfType("*entities.Ticket")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:          "user not found",
			title:         "Request for office supplies",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: 250000,
			requesterID:   uuid.New().String(),
			setupMocks: func() {
				mockUserRepo.On("GetByID", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)
			},
			wantErr:    true,
			errMessage: "user not found",
		},
		{
			name:          "empty title",
			title:         "",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: 250000,
			requesterID:   requesterID,
			setupMocks: func() {
				mockUserRepo.On("GetByID", ctx, requesterID).Return(validUser, nil)
			},
			wantErr:    true,
			errMessage: "title is required",
		},
		{
			name:          "negative cost",
			title:         "Request for office supplies",
			description:   "Need notebooks and pens",
			category:      entities.CategoryOfficeSupplies,
			priority:      entities.PriorityMedium,
			estimatedCost: -1000,
			requesterID:   requesterID,
			setupMocks: func() {
				mockUserRepo.On("GetByID", ctx, requesterID).Return(validUser, nil)
			},
			wantErr:    true,
			errMessage: "estimated cost cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockTicketRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil

			tt.setupMocks()

			ticket, err := service.CreateTicket(
				ctx,
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

				// Verify approval requirements
				expectedApproval := tt.category == entities.CategoryOfficeFurniture || tt.estimatedCost >= 500000
				assert.Equal(t, expectedApproval, ticket.RequiresApproval())

				mockTicketRepo.AssertExpectations(t)
				mockUserRepo.AssertExpectations(t)
			}
		})
	}
}

func TestTicketService_GetTicket(t *testing.T) {
	ctx := context.Background()
	mockTicketRepo := &MockTicketRepository{}
	mockUserRepo := &MockUserRepository{}
	service := services.NewTicketService(mockTicketRepo, mockUserRepo)

	ticketID := uuid.New().String()
	requesterID := uuid.New().String()

	tests := []struct {
		name        string
		ticketID    string
		userID      string
		userRole    entities.UserRole
		setupMocks  func()
		wantErr     bool
		errMessage  string
		expectNil   bool
	}{
		{
			name:     "user viewing own ticket",
			ticketID: ticketID,
			userID:   requesterID,
			userRole: entities.RoleRequester,
			setupMocks: func() {
				ticket := createTestTicket(t, ticketID, requesterID)
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
			},
			wantErr:   false,
			expectNil: false,
		},
		{
			name:     "admin viewing any ticket",
			ticketID: ticketID,
			userID:   uuid.New().String(),
			userRole: entities.RoleAdmin,
			setupMocks: func() {
				ticket := createTestTicket(t, ticketID, requesterID)
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
			},
			wantErr:   false,
			expectNil: false,
		},
		{
			name:     "user viewing someone else's ticket",
			ticketID: ticketID,
			userID:   uuid.New().String(),
			userRole: entities.RoleRequester,
			setupMocks: func() {
				ticket := createTestTicket(t, ticketID, requesterID)
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
			},
			wantErr:   true,
			errMessage: "access denied",
		},
		{
			name:     "ticket not found",
			ticketID: uuid.New().String(),
			userID:   requesterID,
			userRole: entities.RoleRequester,
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)
			},
			wantErr:    true,
			errMessage: "ticket not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockTicketRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil

			tt.setupMocks()

			ticket, err := service.GetTicket(ctx, tt.ticketID, tt.userID, tt.userRole)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				if tt.expectNil {
					assert.Nil(t, ticket)
				} else {
					assert.NotNil(t, ticket)
					assert.Equal(t, tt.ticketID, ticket.GetID())
				}
			}

			if !tt.wantErr || tt.errMessage != "ticket not found" {
				mockTicketRepo.AssertExpectations(t)
			}
		})
	}
}

func TestTicketService_GetUserTickets(t *testing.T) {
	ctx := context.Background()
	mockTicketRepo := &MockTicketRepository{}
	mockUserRepo := &MockUserRepository{}
	service := services.NewTicketService(mockTicketRepo, mockUserRepo)

	requesterID := uuid.New().String()
	tickets := []*entities.Ticket{
		createTestTicket(t, uuid.New().String(), requesterID),
		createTestTicket(t, uuid.New().String(), requesterID),
	}

	tests := []struct {
		name          string
		userID        string
		limit         int
		offset        int
		setupMocks    func()
		wantErr       bool
		expectedCount int
	}{
		{
			name:   "get user tickets successfully",
			userID: requesterID,
			limit:  10,
			offset: 0,
			setupMocks: func() {
				mockTicketRepo.On("GetByRequesterID", ctx, requesterID, 10, 0).Return(tickets, nil)
			},
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name:   "repository error",
			userID: requesterID,
			limit:  10,
			offset: 0,
			setupMocks: func() {
				mockTicketRepo.On("GetByRequesterID", ctx, requesterID, 10, 0).Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockTicketRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil

			tt.setupMocks()

			result, err := service.GetUserTickets(ctx, tt.userID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
				mockTicketRepo.AssertExpectations(t)
			}
		})
	}
}

func TestTicketService_AssignTicket(t *testing.T) {
	ctx := context.Background()
	mockTicketRepo := &MockTicketRepository{}
	mockUserRepo := &MockUserRepository{}
	service := services.NewTicketService(mockTicketRepo, mockUserRepo)

	ticketID := uuid.New().String()
	adminID := uuid.New().String()
	ticket := createTestTicket(t, ticketID, uuid.New().String())
	admin := createTestUser(t, adminID)

	tests := []struct {
		name       string
		ticketID   string
		adminID    string
		setupMocks func()
		wantErr    bool
		errMessage string
	}{
		{
			name:     "successful assignment",
			ticketID: ticketID,
			adminID:  adminID,
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
				mockUserRepo.On("GetByID", ctx, adminID).Return(admin, nil)
				mockTicketRepo.On("Update", ctx, mock.AnythingOfType("*entities.Ticket")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "ticket not found",
			ticketID: uuid.New().String(),
			adminID:  adminID,
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)
			},
			wantErr:    true,
			errMessage: "ticket not found",
		},
		{
			name:     "admin not found",
			ticketID: ticketID,
			adminID:  uuid.New().String(),
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
				mockUserRepo.On("GetByID", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)
			},
			wantErr:    true,
			errMessage: "admin not found",
		},
		{
			name:     "assignee is not admin",
			ticketID: ticketID,
			adminID:  adminID,
			setupMocks: func() {
				regularUser := createTestUser(t, adminID)
				regularUser.ChangeRole(entities.RoleRequester)
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
				mockUserRepo.On("GetByID", ctx, adminID).Return(regularUser, nil)
			},
			wantErr:    true,
			errMessage: "assignee must be an admin",
		},
		{
			name:     "already assigned",
			ticketID: ticketID,
			adminID:  adminID,
			setupMocks: func() {
				assignedTicket := createTestTicket(t, ticketID, uuid.New().String())
				assignedTicket.AssignToAdmin(adminID) // Already assigned
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(assignedTicket, nil)
			},
			wantErr:    true,
			errMessage: "ticket is already assigned",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockTicketRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil

			tt.setupMocks()

			err := service.AssignTicket(ctx, tt.ticketID, tt.adminID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
				mockTicketRepo.AssertExpectations(t)
				mockUserRepo.AssertExpectations(t)
			}
		})
	}
}

func TestTicketService_UpdateTicketStatus(t *testing.T) {
	ctx := context.Background()
	mockTicketRepo := &MockTicketRepository{}
	mockUserRepo := &MockUserRepository{}
	service := services.NewTicketService(mockTicketRepo, mockUserRepo)

	ticketID := uuid.New().String()
	userID := uuid.New().String()
	ticket := createTestTicket(t, ticketID, uuid.New().String())

	tests := []struct {
		name       string
		ticketID   string
		status     entities.TicketStatus
		reason     string
		userID     string
		setupMocks func()
		wantErr    bool
		errMessage string
	}{
		{
			name:     "valid status update",
			ticketID: ticketID,
			status:   entities.StatusInProgress,
			reason:   "Starting work",
			userID:   userID,
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
				mockTicketRepo.On("Update", ctx, mock.AnythingOfType("*entities.Ticket")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "invalid status transition",
			ticketID: ticketID,
			status:   entities.StatusCompleted,
			reason:   "Invalid transition",
			userID:   userID,
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, ticketID).Return(ticket, nil)
			},
			wantErr:    true,
			errMessage: "invalid status transition",
		},
		{
			name:     "ticket not found",
			ticketID: uuid.New().String(),
			status:   entities.StatusInProgress,
			reason:   "Starting work",
			userID:   userID,
			setupMocks: func() {
				mockTicketRepo.On("GetByID", ctx, mock.AnythingOfType("string")).Return(nil, assert.AnError)
			},
			wantErr:    true,
			errMessage: "ticket not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockTicketRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil

			tt.setupMocks()

			err := service.UpdateTicketStatus(ctx, tt.ticketID, tt.status, tt.reason, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
				mockTicketRepo.AssertExpectations(t)
			}
		})
	}
}

// Helper functions
func createTestTicket(t *testing.T, id, requesterID string) *entities.Ticket {
	// Create a simple ticket for testing
	ticket := &entities.Ticket{}
	// In a real implementation, we would create a ticket with the given ID
	// For now, we'll return a mock ticket
	return ticket
}

func createTestUser(t *testing.T, id string) *entities.User {
	// Create a simple user for testing
	user := &entities.User{}
	// In a real implementation, we would create a user with the given ID
	// For now, we'll return a mock user
	return user
}