package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/infrastructure/auth"
)

func TestUser_NewUser(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)

	tests := []struct {
		name         string
		employeeID   string
		userName     string
		email        string
		department   string
		role         entities.UserRole
		password     string
		wantErr      bool
		errMessage   string
	}{
		{
			name:       "valid requester",
			employeeID: "EMP001",
			userName:   "John Doe",
			email:      "john.doe@company.com",
			department: "Engineering",
			role:       entities.RoleRequester,
			password:   "SecurePass123!",
			wantErr:    false,
		},
		{
			name:       "valid admin",
			employeeID: "ADMIN001",
			userName:   "Admin User",
			email:      "admin@company.com",
			department: "IT",
			role:       entities.RoleAdmin,
			password:   "AdminPass123!",
			wantErr:    false,
		},
		{
			name:       "valid approver",
			employeeID: "APP001",
			userName:   "Approver User",
			email:      "approver@company.com",
			department: "Finance",
			role:       entities.RoleApprover,
			password:   "ApproverPass123!",
			wantErr:    false,
		},
		{
			name:       "empty employee ID",
			employeeID: "",
			userName:   "John Doe",
			email:      "john.doe@company.com",
			department: "Engineering",
			role:       entities.RoleRequester,
			password:   "SecurePass123!",
			wantErr:    true,
			errMessage: "employee ID is required",
		},
		{
			name:       "empty name",
			employeeID: "EMP001",
			userName:   "",
			email:      "john.doe@company.com",
			department: "Engineering",
			role:       entities.RoleRequester,
			password:   "SecurePass123!",
			wantErr:    true,
			errMessage: "name is required",
		},
		{
			name:       "invalid email",
			employeeID: "EMP001",
			userName:   "John Doe",
			email:      "invalid-email",
			department: "Engineering",
			role:       entities.RoleRequester,
			password:   "SecurePass123!",
			wantErr:    true,
			errMessage: "invalid email format",
		},
		{
			name:       "weak password",
			employeeID: "EMP001",
			userName:   "John Doe",
			email:      "john.doe@company.com",
			department: "Engineering",
			role:       entities.RoleRequester,
			password:   "123",
			wantErr:    true,
			errMessage: "password does not meet security requirements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entities.NewUser(
				tt.employeeID,
				tt.userName,
				tt.email,
				tt.department,
				tt.role,
				tt.password,
				passwordHasher,
			)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMessage)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.employeeID, user.GetEmployeeID())
				assert.Equal(t, tt.userName, user.GetName())
				assert.Equal(t, tt.email, user.GetEmail())
				assert.Equal(t, tt.department, user.GetDepartment())
				assert.Equal(t, tt.role, user.GetRole())
				assert.True(t, user.IsActive())
				assert.NotEmpty(t, user.GetID())
				assert.NotEmpty(t, user.GetPasswordHash())
				assert.False(t, user.GetPasswordHash() == tt.password) // Should be hashed
			}
		})
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	password := "SecurePass123!"

	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		password,
		passwordHasher,
	)
	require.NoError(t, err)

	// Correct password should verify
	assert.True(t, user.VerifyPassword(password, passwordHasher))

	// Incorrect password should not verify
	assert.False(t, user.VerifyPassword("WrongPassword", passwordHasher))
}

func TestUser_ChangePassword(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	oldPassword := "SecurePass123!"
	newPassword := "NewSecurePass456!"

	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		oldPassword,
		passwordHasher,
	)
	require.NoError(t, err)

	// Test password change with correct old password
	err = user.ChangePassword(oldPassword, newPassword, passwordHasher)
	assert.NoError(t, err)

	// Verify old password no longer works
	assert.False(t, user.VerifyPassword(oldPassword, passwordHasher))

	// Verify new password works
	assert.True(t, user.VerifyPassword(newPassword, passwordHasher))

	// Test password change with incorrect old password
	err = user.ChangePassword("WrongPassword", "AnotherNewPass", passwordHasher)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "current password is incorrect")

	// Test password change with weak new password
	err = user.ChangePassword(newPassword, "weak", passwordHasher)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "new password does not meet security requirements")
}

func TestUser_UpdateProfile(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		"SecurePass123!",
		passwordHasher,
	)
	require.NoError(t, err)

	// Test valid profile update
	err = user.UpdateProfile("John Smith", "john.smith@company.com", "Product")
	assert.NoError(t, err)
	assert.Equal(t, "John Smith", user.GetName())
	assert.Equal(t, "john.smith@company.com", user.GetEmail())
	assert.Equal(t, "Product", user.GetDepartment())

	// Test empty name
	err = user.UpdateProfile("", "john.smith@company.com", "Product")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test invalid email
	err = user.UpdateProfile("John Smith", "invalid-email", "Product")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format")
}

func TestUser_ChangeRole(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		"SecurePass123!",
		passwordHasher,
	)
	require.NoError(t, err)

	// Test role change
	err = user.ChangeRole(entities.RoleAdmin)
	assert.NoError(t, err)
	assert.Equal(t, entities.RoleAdmin, user.GetRole())

	// Test role change to same role
	err = user.ChangeRole(entities.RoleAdmin)
	assert.NoError(t, err)
	assert.Equal(t, entities.RoleAdmin, user.GetRole())
}

func TestUser_Deactivate(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		"SecurePass123!",
		passwordHasher,
	)
	require.NoError(t, err)

	// Test deactivation
	assert.True(t, user.IsActive())
	user.Deactivate()
	assert.False(t, user.IsActive())
}

func TestUser_Activate(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		"SecurePass123!",
		passwordHasher,
	)
	require.NoError(t, err)

	// Deactivate first
	user.Deactivate()
	assert.False(t, user.IsActive())

	// Test activation
	user.Activate()
	assert.True(t, user.IsActive())
}

func TestUser_HasPermission(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)

	tests := []struct {
		name      string
		role      entities.UserRole
		permission string
		expected  bool
	}{
		{
			name:       "requester can create tickets",
			role:       entities.RoleRequester,
			permission: "create_ticket",
			expected:   true,
		},
		{
			name:       "requester cannot assign tickets",
			role:       entities.RoleRequester,
			permission: "assign_ticket",
			expected:   false,
		},
		{
			name:       "requester cannot approve tickets",
			role:       entities.RoleRequester,
			permission: "approve_ticket",
			expected:   false,
		},
		{
			name:       "admin can assign tickets",
			role:       entities.RoleAdmin,
			permission: "assign_ticket",
			expected:   true,
		},
		{
			name:       "admin can approve tickets",
			role:       entities.RoleAdmin,
			permission: "approve_ticket",
			expected:   true,
		},
		{
			name:       "admin can manage assets",
			role:       entities.RoleAdmin,
			permission: "manage_assets",
			expected:   true,
		},
		{
			name:       "approver can approve tickets",
			role:       entities.RoleApprover,
			permission: "approve_ticket",
			expected:   true,
		},
		{
			name:       "approver cannot manage assets",
			role:       entities.RoleApprover,
			permission: "manage_assets",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entities.NewUser(
				"EMP001",
				"Test User",
				"test@company.com",
				"Engineering",
				tt.role,
				"SecurePass123!",
				passwordHasher,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, user.HasPermission(tt.permission))
		})
	}
}

func TestUser_CanViewTicket(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	userID := uuid.New()
	ticketRequesterID := uuid.New()

	// Create different user types
	requester, _ := entities.NewUser(
		"EMP001",
		"Requester",
		"requester@company.com",
		"Engineering",
		entities.RoleRequester,
		"SecurePass123!",
		passwordHasher,
	)

	admin, _ := entities.NewUser(
		"ADMIN001",
		"Admin",
		"admin@company.com",
		"IT",
		entities.RoleAdmin,
		"AdminPass123!",
		passwordHasher,
	)

	approver, _ := entities.NewUser(
		"APP001",
		"Approver",
		"approver@company.com",
		"Finance",
		entities.RoleApprover,
		"ApproverPass123!",
		passwordHasher,
	)

	// Test viewing own ticket
	assert.True(t, requester.CanViewTicket(userID.String(), ticketRequesterID.String()))

	// Test viewing someone else's ticket as requester
	assert.False(t, requester.CanViewTicket(uuid.New().String(), ticketRequesterID.String()))

	// Test admin can view any ticket
	assert.True(t, admin.CanViewTicket(uuid.New().String(), ticketRequesterID.String()))

	// Test approver can view tickets requiring approval (simplified test)
	assert.True(t, approver.CanViewTicket(uuid.New().String(), ticketRequesterID.String()))
}

func TestUser_GetUserInfo(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)
	user, err := entities.NewUser(
		"EMP001",
		"John Doe",
		"john.doe@company.com",
		"Engineering",
		entities.RoleRequester,
		"SecurePass123!",
		passwordHasher,
	)
	require.NoError(t, err)

	userInfo := user.GetUserInfo()
	assert.Equal(t, user.GetID(), userInfo.ID)
	assert.Equal(t, user.GetEmployeeID(), userInfo.EmployeeID)
	assert.Equal(t, user.GetName(), userInfo.Name)
	assert.Equal(t, user.GetEmail(), userInfo.Email)
	assert.Equal(t, string(user.GetRole()), userInfo.Role)
	assert.Equal(t, user.GetDepartment(), userInfo.Department)
}

func TestUser_IsRole(t *testing.T) {
	passwordHasher := auth.NewPasswordHasher(nil)

	tests := []struct {
		name     string
		role     entities.UserRole
		checkRole entities.UserRole
		expected bool
	}{
		{"requester is requester", entities.RoleRequester, entities.RoleRequester, true},
		{"requester is not admin", entities.RoleRequester, entities.RoleAdmin, false},
		{"admin is admin", entities.RoleAdmin, entities.RoleAdmin, true},
		{"admin is not approver", entities.RoleAdmin, entities.RoleApprover, false},
		{"approver is approver", entities.RoleApprover, entities.RoleApprover, true},
		{"approver is not requester", entities.RoleApprover, entities.RoleRequester, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entities.NewUser(
				"EMP001",
				"Test User",
				"test@company.com",
				"Engineering",
				tt.role,
				"SecurePass123!",
				passwordHasher,
			)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, user.IsRole(tt.checkRole))
		})
	}
}