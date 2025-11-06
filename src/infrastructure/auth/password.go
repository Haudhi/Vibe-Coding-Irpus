package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
	"go.uber.org/zap"
)

// PasswordHasher handles password hashing and verification
type PasswordHasher struct {
	logger *zap.Logger
}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher(logger *zap.Logger) *PasswordHasher {
	return &PasswordHasher{
		logger: logger,
	}
}

// HashConfig contains configuration for password hashing
type HashConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

// DefaultHashConfig returns secure default hashing configuration
func DefaultHashConfig() *HashConfig {
	return &HashConfig{
		Time:    3,       // Number of iterations
		Memory:  64 * 1024, // 64MB
		Threads: 4,       // Number of threads
		KeyLen:  32,      // 32 bytes
		SaltLen: 16,      // 16 bytes
	}
}

// HashPassword hashes a password using Argon2id
func (p *PasswordHasher) HashPassword(password string) (string, error) {
	config := DefaultHashConfig()

	// Generate random salt
	salt := make([]byte, config.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		p.logger.Error("Failed to generate salt for password hashing", zap.Error(err))
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password
	hash := argon2.IDKey([]byte(password), salt, config.Time, config.Memory, config.Threads, config.KeyLen)

	// Encode salt and hash
	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	// Combine them into a stored password format
	storedPassword := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.Memory, config.Time, config.Threads, saltBase64, hashBase64)

	p.logger.Debug("Password hashed successfully")
	return storedPassword, nil
}

// VerifyPassword verifies a password against its hash
func (p *PasswordHasher) VerifyPassword(password, hashedPassword string) (bool, error) {
	// Parse the hashed password
	var version int
	var memory, time uint32
	var threads uint8
	var salt, hash []byte

	_, err := fmt.Sscanf(hashedPassword, "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		&version, &memory, &time, &threads, &salt, &hash)
	if err != nil {
		p.logger.Warn("Failed to parse hashed password format", zap.Error(err))
		return false, fmt.Errorf("invalid hashed password format: %w", err)
	}

	// Decode salt and hash
	decodedSalt, err := base64.RawStdEncoding.DecodeString(string(salt))
	if err != nil {
		p.logger.Warn("Failed to decode salt from hashed password", zap.Error(err))
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(string(hash))
	if err != nil {
		p.logger.Warn("Failed to decode hash from hashed password", zap.Error(err))
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the provided password with the same salt
	hashedInput := argon2.IDKey([]byte(password), decodedSalt, time, memory, threads, uint32(len(decodedHash)))

	// Compare the hashes using constant-time comparison
	match := subtle.ConstantTimeCompare(hashedInput, decodedHash) == 1

	if match {
		p.logger.Debug("Password verification successful")
	} else {
		p.logger.Debug("Password verification failed")
	}

	return match, nil
}

// GenerateRandomToken generates a cryptographically secure random token
func (p *PasswordHasher) GenerateRandomToken(length int) (string, error) {
	if length <= 0 {
		length = 32 // Default length
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		p.logger.Error("Failed to generate random token", zap.Error(err))
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	p.logger.Debug("Random token generated successfully")
	return token, nil
}

// ValidatePasswordStrength checks password strength requirements
func (p *PasswordHasher) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be no more than 128 characters long")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	// Check for at least one lowercase letter
	hasLower := false
	// Check for at least one digit
	hasDigit := false
	// Check for at least one special character
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= ' ' && char <= '~':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}

	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	p.logger.Debug("Password strength validation passed")
	return nil
}

// PasswordStrength represents password strength levels
type PasswordStrength int

const (
	PasswordStrengthWeak PasswordStrength = iota
	PasswordStrengthFair
	PasswordStrengthGood
	PasswordStrengthStrong
)

// CalculatePasswordStrength calculates password strength
func (p *PasswordHasher) CalculatePasswordStrength(password string) PasswordStrength {
	strength := 0

	// Length factor
	if len(password) >= 8 {
		strength++
	}
	if len(password) >= 12 {
		strength++
	}
	if len(password) >= 16 {
		strength++
	}

	// Character variety factor
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= ' ' && char <= '~':
			hasSpecial = true
		}
	}

	if hasUpper {
		strength++
	}
	if hasLower {
		strength++
	}
	if hasDigit {
		strength++
	}
	if hasSpecial {
		strength++
	}

	// Convert strength score to level
	switch {
	case strength >= 7:
		return PasswordStrengthStrong
	case strength >= 5:
		return PasswordStrengthGood
	case strength >= 3:
		return PasswordStrengthFair
	default:
		return PasswordStrengthWeak
	}
}

// GetPasswordStrengthMessage returns a human-readable password strength message
func (p *PasswordHasher) GetPasswordStrengthMessage(strength PasswordStrength) string {
	switch strength {
	case PasswordStrengthWeak:
		return "Weak password - consider adding more characters and variety"
	case PasswordStrengthFair:
		return "Fair password - could be stronger"
	case PasswordStrengthGood:
		return "Good password"
	case PasswordStrengthStrong:
		return "Strong password"
	default:
		return "Unknown password strength"
	}
}