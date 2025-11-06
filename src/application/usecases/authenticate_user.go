package usecases

import (
	"fmt"
	"time"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/repositories"
	"github.com/company/ga-ticketing/src/infrastructure/auth"
	"go.uber.org/zap"
)

// AuthenticateUserUseCase handles user authentication
type AuthenticateUserUseCase struct {
	userRepo       repositories.UserRepository
	jwtManager     *auth.JWTManager
	passwordHasher *auth.PasswordHasher
	logger         *zap.Logger
}

// NewAuthenticateUserUseCase creates a new AuthenticateUserUseCase
func NewAuthenticateUserUseCase(
	userRepo repositories.UserRepository,
	jwtManager *auth.JWTManager,
	passwordHasher *auth.PasswordHasher,
	logger *zap.Logger,
) *AuthenticateUserUseCase {
	return &AuthenticateUserUseCase{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		passwordHasher: passwordHasher,
		logger:         logger,
	}
}

// Execute authenticates a user and returns a login response
func (uc *AuthenticateUserUseCase) Execute(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Validate input
	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Find user by email
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		uc.logger.Warn("Login attempt with non-existent email",
			zap.String("email", req.Email),
			zap.Error(err))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive() {
		uc.logger.Warn("Login attempt with inactive user",
			zap.String("user_id", user.GetID()),
			zap.String("email", req.Email))
		return nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	if !user.VerifyPassword(req.Password, uc.passwordHasher) {
		uc.logger.Warn("Login attempt with invalid password",
			zap.String("user_id", user.GetID()),
			zap.String("email", req.Email))
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	userInfo := user.GetUserInfo()
	token, err := uc.jwtManager.GenerateToken(userInfo)
	if err != nil {
		uc.logger.Error("Failed to generate JWT token during login",
			zap.String("user_id", user.GetID()),
			zap.Error(err))
		return nil, fmt.Errorf("failed to generate authentication token")
	}

	uc.logger.Info("User authenticated successfully",
		zap.String("user_id", user.GetID()),
		zap.String("email", req.Email),
		zap.String("role", string(user.GetRole())))

	// Prepare response
	response := &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.GetID(),
			Email:     user.GetEmail(),
			Name:      user.GetName(),
			Role:      string(user.GetRole()),
			Department: user.GetDepartment(),
			CreatedAt: time.Now(), // This should come from the actual created_at field in database
		},
	}

	return response, nil
}