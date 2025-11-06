package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/company/ga-ticketing/src/infrastructure/auth"
)

// User represents the authenticated user in the context
type User struct {
	ID    string
	Email string
	Name  string
	Role  string
}

// userContextKey is the key used to store user in context
type userContextKey string

const userKey userContextKey = "user"

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
	logger     *zap.Logger
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtManager *auth.JWTManager, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
		logger:     logger,
	}
}

// Authenticate middleware that validates JWT tokens
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Authorization header required"})
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString, err := m.jwtManager.ExtractTokenFromHeader(authHeader)
		if err != nil {
			m.logger.Warn("Invalid authorization header format", zap.Error(err))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Invalid authorization header format"})
			return
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			m.logger.Warn("Invalid JWT token", zap.Error(err))
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": "Invalid or expired token"})
			return
		}

		// Create user struct and add to context
		user := &User{
			ID:    claims.UserID,
			Email: claims.Email,
			Name:  claims.Name,
			Role:  claims.Role,
		}

		ctx := context.WithValue(r.Context(), userKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware that checks if user has required role
func (m *AuthMiddleware) RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("user_role").(string)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "User role not found"})
				return
			}

			if userRole != requiredRole {
				m.logger.Warn("Insufficient permissions",
					zap.String("user_role", userRole),
					zap.String("required_role", requiredRole),
				)
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, map[string]string{"error": "Insufficient permissions"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireOneOfRoles middleware that checks if user has any of the required roles
func (m *AuthMiddleware) RequireOneOfRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("user_role").(string)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "User role not found"})
				return
			}

			hasRequiredRole := false
			for _, role := range roles {
				if userRole == role {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				m.logger.Warn("Insufficient permissions",
					zap.String("user_role", userRole),
					zap.Strings("required_roles", roles),
				)
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, map[string]string{"error": "Insufficient permissions"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission middleware that checks if user has specific permission
func (m *AuthMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("user_role").(string)
			if !ok {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "User role not found"})
				return
			}

			if !m.hasPermission(userRole, permission) {
				m.logger.Warn("Insufficient permissions",
					zap.String("user_role", userRole),
					zap.String("required_permission", permission),
				)
				render.Status(r, http.StatusForbidden)
				render.JSON(w, r, map[string]string{"error": "Insufficient permissions"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuth middleware that doesn't require authentication but adds user info if token is present
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header, proceed without user info
			next.ServeHTTP(w, r)
			return
		}

		// Try to extract and validate token
		tokenString, err := m.jwtManager.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Invalid format, proceed without user info
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			// Invalid token, proceed without user info
			next.ServeHTTP(w, r)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_role", claims.Role)
		ctx = context.WithValue(ctx, "user_email", claims.Email)
		ctx = context.WithValue(ctx, "user_name", claims.Name)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CORS middleware for handling cross-origin requests
func (m *AuthMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// hasPermission checks if a role has a specific permission
func (m *AuthMiddleware) hasPermission(userRole, permission string) bool {
	permissions := getRolePermissions(userRole)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// getRolePermissions returns permissions for each role
func getRolePermissions(role string) []string {
	switch role {
	case "requester":
		return []string{
			"create_ticket",
			"view_own_tickets",
			"add_comments_to_own_tickets",
		}
	case "approver":
		return []string{
			"create_ticket",
			"view_own_tickets",
			"add_comments_to_own_tickets",
			"view_tickets_for_approval",
			"approve_ticket",
			"reject_ticket",
		}
	case "admin":
		return []string{
			"create_ticket",
			"view_own_tickets",
			"add_comments_to_own_tickets",
			"view_all_tickets",
			"assign_ticket",
			"update_ticket_status",
			"approve_ticket",
			"reject_ticket",
			"manage_assets",
			"view_all_users",
			"manage_users",
		}
	default:
		return []string{}
	}
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(userKey).(*User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}