package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	middleware "github.com/company/ga-ticketing/src/interface/http/middleware"
	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/application/usecases"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authenticateUserUC *usecases.AuthenticateUserUseCase
	getCurrentUserUC   *usecases.GetCurrentUserUseCase
	logger             *zap.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	authenticateUserUC *usecases.AuthenticateUserUseCase,
	getCurrentUserUC *usecases.GetCurrentUserUseCase,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		authenticateUserUC: authenticateUserUC,
		getCurrentUserUC:   getCurrentUserUC,
		logger:             logger,
	}
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Error("Failed to decode login request", zap.Error(err))
		render.Render(w, r, middleware.ErrBadRequest(err.Error()))
		return
	}

	// Validate request
	if req.Email == "" {
		render.Render(w, r, middleware.ErrBadRequest("email is required"))
		return
	}
	if req.Password == "" {
		render.Render(w, r, middleware.ErrBadRequest("password is required"))
		return
	}

	// Authenticate user
	response, err := h.authenticateUserUC.Execute(req)
	if err != nil {
		h.logger.Warn("Authentication failed",
			zap.String("email", req.Email),
			zap.Error(err))

		// Don't reveal specific error for security reasons
		render.Render(w, r, middleware.ErrUnauthorized("Invalid credentials"))
		return
	}

	h.logger.Info("User logged in successfully",
		zap.String("email", req.Email),
		zap.String("user_id", response.User.ID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// GetCurrentUser handles getting current user information
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		h.logger.Error("Failed to get user from context", zap.Error(err))
		render.Render(w, r, middleware.ErrUnauthorized("User not found in context"))
		return
	}

	// Get current user information
	response, err := h.getCurrentUserUC.Execute(user.ID)
	if err != nil {
		h.logger.Error("Failed to get current user information",
			zap.String("user_id", user.ID),
			zap.Error(err))
		render.Render(w, r, middleware.ErrNotFound("User not found"))
		return
	}

	h.logger.Debug("Retrieved current user information",
		zap.String("user_id", user.ID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}