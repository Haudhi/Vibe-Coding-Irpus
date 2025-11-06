package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Config holds JWT configuration
type Config struct {
	Secret         string
	Expiry         time.Duration
	RefreshExpiry  time.Duration
	Issuer         string
}

// DefaultConfig returns default JWT configuration
func DefaultConfig() *Config {
	return &Config{
		Secret:        "your-super-secret-jwt-key-change-this-in-production",
		Expiry:        24 * time.Hour,
		RefreshExpiry: 7 * 24 * time.Hour,
		Issuer:        "ga-ticketing-system",
	}
}

// Claims represents JWT claims structure
type Claims struct {
	UserID     string `json:"user_id"`
	EmployeeID string `json:"employee_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	jwt.RegisteredClaims
}

// UserInfo represents user information for token generation
type UserInfo struct {
	ID         string
	EmployeeID string
	Name       string
	Email      string
	Role       string
	Department string
}

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	config *Config
	logger *zap.Logger
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(config *Config, logger *zap.Logger) *JWTManager {
	if config == nil {
		config = DefaultConfig()
	}

	return &JWTManager{
		config: config,
		logger: logger,
	}
}

// GenerateToken generates a new JWT access token
func (j *JWTManager) GenerateToken(userInfo UserInfo) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.config.Expiry)

	claims := &Claims{
		UserID:     userInfo.ID,
		EmployeeID: userInfo.EmployeeID,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Role:       userInfo.Role,
		Department: userInfo.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.config.Issuer,
			Subject:   userInfo.ID,
			Audience:  []string{"ga-ticketing-client"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		j.logger.Error("Failed to generate JWT token",
			zap.String("user_id", userInfo.ID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	j.logger.Info("JWT token generated successfully",
		zap.String("user_id", userInfo.ID),
		zap.String("token_id", claims.ID),
		zap.Time("expires_at", expiresAt),
	)

	return tokenString, nil
}

// GenerateRefreshToken generates a new JWT refresh token
func (j *JWTManager) GenerateRefreshToken(userInfo UserInfo) (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.config.RefreshExpiry)

	claims := &Claims{
		UserID:     userInfo.ID,
		EmployeeID: userInfo.EmployeeID,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Role:       userInfo.Role,
		Department: userInfo.Department,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    j.config.Issuer,
			Subject:   userInfo.ID,
			Audience:  []string{"ga-ticketing-refresh"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.config.Secret))
	if err != nil {
		j.logger.Error("Failed to generate JWT refresh token",
			zap.String("user_id", userInfo.ID),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	j.logger.Info("JWT refresh token generated successfully",
		zap.String("user_id", userInfo.ID),
		zap.String("token_id", claims.ID),
		zap.Time("expires_at", expiresAt),
	)

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Secret), nil
	})

	if err != nil {
		j.logger.Warn("Failed to parse JWT token", zap.Error(err))
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Additional validations
	if err := j.validateClaims(claims); err != nil {
		return nil, err
	}

	j.logger.Debug("JWT token validated successfully",
		zap.String("user_id", claims.UserID),
		zap.String("token_id", claims.ID),
		zap.Time("expires_at", claims.ExpiresAt.Time),
	)

	return claims, nil
}

// validateClaims performs additional claim validations
func (j *JWTManager) validateClaims(claims *Claims) error {
	now := time.Now()

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(now) {
		return fmt.Errorf("token has expired")
	}

	// Check if token is not yet valid
	if claims.NotBefore != nil && claims.NotBefore.After(now) {
		return fmt.Errorf("token is not yet valid")
	}

	// Check issuer
	if claims.Issuer != j.config.Issuer {
		return fmt.Errorf("invalid token issuer")
	}

	// Validate required fields
	if claims.UserID == "" {
		return fmt.Errorf("missing user ID in token")
	}

	if claims.Email == "" {
		return fmt.Errorf("missing email in token")
	}

	if claims.Role == "" {
		return fmt.Errorf("missing role in token")
	}

	// Validate role
	validRoles := map[string]bool{
		"requester": true,
		"approver":  true,
		"admin":     true,
	}
	if !validRoles[claims.Role] {
		return fmt.Errorf("invalid role: %s", claims.Role)
	}

	return nil
}

// RefreshToken generates a new access token from a refresh token
func (j *JWTManager) RefreshToken(refreshTokenString string) (string, error) {
	claims, err := j.ValidateToken(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if this is a refresh token
	if !contains(claims.Audience, "ga-ticketing-refresh") {
		return "", fmt.Errorf("not a refresh token")
	}

	// Generate new access token
	userInfo := UserInfo{
		ID:         claims.UserID,
		EmployeeID: claims.EmployeeID,
		Name:       claims.Name,
		Email:      claims.Email,
		Role:       claims.Role,
		Department: claims.Department,
	}

	return j.GenerateToken(userInfo)
}

// GetTokenExpiration returns the expiration time of a token
func (j *JWTManager) GetTokenExpiration(tokenString string) (*time.Time, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt != nil {
		expiresAt := claims.ExpiresAt.Time
		return &expiresAt, nil
	}

	return nil, fmt.Errorf("no expiration claim found")
}

// IsTokenExpired checks if a token is expired
func (j *JWTManager) IsTokenExpired(tokenString string) bool {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return true // Consider invalid tokens as expired
	}

	if claims.ExpiresAt == nil {
		return true
	}

	return claims.ExpiresAt.Before(time.Now())
}

// ExtractTokenFromHeader extracts token from Authorization header
func (j *JWTManager) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}