package entities

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/company/ga-ticketing/src/infrastructure/auth"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleRequester UserRole = "requester"
	RoleApprover  UserRole = "approver"
	RoleAdmin     UserRole = "admin"
)

// User represents a system user
type User struct {
	id           string
	employeeID   string
	name         string
	email        string
	department   string
	role         UserRole
	passwordHash string
	isActive     bool
}

// NewUser creates a new user with password hashing
func NewUser(employeeID, name, email, department string, role UserRole, password string, passwordHasher *auth.PasswordHasher) (*User, error) {
	// Validate input
	if employeeID == "" {
		return nil, errors.New("employee ID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}
	if !isValidRole(role) {
		return nil, errors.New("invalid role")
	}

	// Validate password strength
	if err := passwordHasher.ValidatePasswordStrength(password); err != nil {
		return nil, fmt.Errorf("password does not meet security requirements: %w", err)
	}

	// Hash password
	hashedPassword, err := passwordHasher.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return &User{
		id:           generateUserID(),
		employeeID:   employeeID,
		name:         name,
		email:        email,
		department:   department,
		role:         role,
		passwordHash: hashedPassword,
		isActive:     true,
	}, nil
}

// Getters
func (u *User) GetID() string           { return u.id }
func (u *User) GetEmployeeID() string   { return u.employeeID }
func (u *User) GetName() string         { return u.name }
func (u *User) GetEmail() string        { return u.email }
func (u *User) GetDepartment() string   { return u.department }
func (u *User) GetRole() UserRole       { return u.role }
func (u *User) GetPasswordHash() string { return u.passwordHash }
func (u *User) IsActive() bool          { return u.isActive }

// VerifyPassword checks if the provided password matches the user's password
func (u *User) VerifyPassword(password string, passwordHasher *auth.PasswordHasher) bool {
	isValid, err := passwordHasher.VerifyPassword(password, u.passwordHash)
	if err != nil {
		return false
	}
	return isValid
}

// ChangePassword updates the user's password
func (u *User) ChangePassword(currentPassword, newPassword string, passwordHasher *auth.PasswordHasher) error {
	// Verify current password
	if !u.VerifyPassword(currentPassword, passwordHasher) {
		return errors.New("current password is incorrect")
	}

	// Validate new password
	if err := passwordHasher.ValidatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("new password does not meet security requirements: %w", err)
	}

	// Hash new password
	newHashedPassword, err := passwordHasher.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	u.passwordHash = newHashedPassword
	return nil
}

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(name, email, department string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}

	u.name = name
	u.email = email
	u.department = department
	return nil
}

// ChangeRole updates the user's role
func (u *User) ChangeRole(newRole UserRole) error {
	if !isValidRole(newRole) {
		return errors.New("invalid role")
	}
	u.role = newRole
	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.isActive = false
}

// Activate activates the user account
func (u *User) Activate() {
	u.isActive = true
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permission string) bool {
	permissions := getRolePermissions(u.role)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// CanViewTicket checks if the user can view a specific ticket
func (u *User) CanViewTicket(userID, ticketRequesterID string) bool {
	// Admin can view any ticket
	if u.role == RoleAdmin {
		return true
	}

	// User can view their own tickets
	if u.id == userID && u.id == ticketRequesterID {
		return true
	}

	// Approver can view tickets requiring approval
	if u.role == RoleApprover {
		// In a real implementation, we would check if the ticket requires approval
		return true
	}

	return false
}

// IsRole checks if the user has a specific role
func (u *User) IsRole(role UserRole) bool {
	return u.role == role
}

// GetUserInfo returns user information suitable for JWT claims
func (u *User) GetUserInfo() auth.UserInfo {
	return auth.UserInfo{
		ID:         u.id,
		EmployeeID: u.employeeID,
		Name:       u.name,
		Email:      u.email,
		Role:       string(u.role),
		Department: u.department,
	}
}

// Helper functions
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidRole(role UserRole) bool {
	switch role {
	case RoleRequester, RoleApprover, RoleAdmin:
		return true
	default:
		return false
	}
}

func generateUserID() string {
	return uuid.New().String()
}

// getRolePermissions returns the permissions for each role
func getRolePermissions(role UserRole) []string {
	switch role {
	case RoleRequester:
		return []string{
			"create_ticket",
			"view_own_tickets",
			"add_comments_to_own_tickets",
			"view_profile",
			"update_profile",
		}
	case RoleApprover:
		return []string{
			"create_ticket",
			"view_own_tickets",
			"add_comments_to_own_tickets",
			"view_profile",
			"update_profile",
			"view_tickets_for_approval",
			"approve_ticket",
			"reject_ticket",
		}
	case RoleAdmin:
		return []string{
			"create_ticket",
			"view_own_tickets",
			"add_comments_to_own_tickets",
			"view_profile",
			"update_profile",
			"view_all_tickets",
			"assign_ticket",
			"update_ticket_status",
			"approve_ticket",
			"reject_ticket",
			"manage_assets",
			"view_all_users",
			"manage_users",
			"view_reports",
		}
	default:
		return []string{}
	}
}

// GetAllRoles returns all available user roles
func GetAllRoles() []UserRole {
	return []UserRole{
		RoleRequester,
		RoleApprover,
		RoleAdmin,
	}
}

// RoleFromString converts a string to UserRole
func RoleFromString(roleStr string) (UserRole, error) {
	roleStr = strings.ToLower(strings.TrimSpace(roleStr))
	switch roleStr {
	case "requester":
		return RoleRequester, nil
	case "approver":
		return RoleApprover, nil
	case "admin":
		return RoleAdmin, nil
	default:
		return RoleRequester, fmt.Errorf("invalid role: %s", roleStr)
	}
}

// RoleDisplayName returns a human-readable display name for a role
func RoleDisplayName(role UserRole) string {
	switch role {
	case RoleRequester:
		return "Requester"
	case RoleApprover:
		return "Approver"
	case RoleAdmin:
		return "Administrator"
	default:
		return string(role)
	}
}