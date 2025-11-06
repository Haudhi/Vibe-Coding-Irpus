package usecases

import (
	"fmt"
	"time"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/repositories"
	"go.uber.org/zap"
)

// GetCurrentUserUseCase handles retrieving current user information
type GetCurrentUserUseCase struct {
	userRepo repositories.UserRepository
	logger   *zap.Logger
}

// NewGetCurrentUserUseCase creates a new GetCurrentUserUseCase
func NewGetCurrentUserUseCase(
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// Execute retrieves the current user information by user ID
func (uc *GetCurrentUserUseCase) Execute(userID string) (*dto.UserResponse, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Find user by ID
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		uc.logger.Error("Failed to find user for get current user request",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive() {
		uc.logger.Warn("Attempt to get current user for inactive account",
			zap.String("user_id", userID))
		return nil, fmt.Errorf("account is inactive")
	}

	uc.logger.Debug("Retrieved current user information",
		zap.String("user_id", userID),
		zap.String("email", user.GetEmail()))

	// Prepare response
	response := &dto.UserResponse{
		ID:        user.GetID(),
		Email:     user.GetEmail(),
		Name:      user.GetName(),
		Role:      string(user.GetRole()),
		Department: user.GetDepartment(),
		CreatedAt: time.Now(), // This should come from the actual created_at field in database
	}

	return response, nil
}