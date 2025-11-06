package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/repositories"
	"github.com/company/ga-ticketing/src/infrastructure/auth"
)

// PGUserRepository implements the UserRepository interface for PostgreSQL
type PGUserRepository struct {
	db *pgxpool.Pool
}

// NewPGUserRepository creates a new PostgreSQL user repository
func NewPGUserRepository(db *pgxpool.Pool) repositories.UserRepository {
	return &PGUserRepository{db: db}
}

// Create saves a new user to the database
func (r *PGUserRepository) Create(user *entities.User) error {
	query := `
		INSERT INTO users (id, employee_id, name, email, department, role, password_hash, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(context.Background(), query,
		user.GetID(),
		user.GetEmployeeID(),
		user.GetName(),
		user.GetEmail(),
		user.GetDepartment(),
		user.GetRole(),
		user.GetPasswordHash(),
		user.IsActive(),
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by their ID
func (r *PGUserRepository) FindByID(id string) (*entities.User, error) {
	query := `
		SELECT id, employee_id, name, email, department, role, password_hash, is_active
		FROM users
		WHERE id = $1 AND is_active = true
	`

	var user entities.User
	var passwordHash sql.NullString

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.EmployeeID,
		&user.Name,
		&user.Email,
		&user.Department,
		&user.Role,
		&passwordHash,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
	}

	return &user, nil
}

// FindByEmail retrieves a user by their email
func (r *PGUserRepository) FindByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id, employee_id, name, email, department, role, password_hash, is_active
		FROM users
		WHERE email = $1 AND is_active = true
	`

	var user entities.User
	var passwordHash sql.NullString

	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.EmployeeID,
		&user.Name,
		&user.Email,
		&user.Department,
		&user.Role,
		&passwordHash,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
	}

	return &user, nil
}

// FindByEmployeeID retrieves a user by their employee ID
func (r *PGUserRepository) FindByEmployeeID(employeeID string) (*entities.User, error) {
	query := `
		SELECT id, employee_id, name, email, department, role, password_hash, is_active
		FROM users
		WHERE employee_id = $1 AND is_active = true
	`

	var user entities.User
	var passwordHash sql.NullString

	err := r.db.QueryRow(context.Background(), query, employeeID).Scan(
		&user.ID,
		&user.EmployeeID,
		&user.Name,
		&user.Email,
		&user.Department,
		&user.Role,
		&passwordHash,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if passwordHash.Valid {
		user.PasswordHash = passwordHash.String
	}

	return &user, nil
}

// FindAll retrieves all users with optional filtering
func (r *PGUserRepository) FindAll(role *entities.UserRole, isActive *bool) ([]*entities.User, error) {
	query := `
		SELECT id, employee_id, name, email, department, role, password_hash, is_active
		FROM users
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if role != nil {
		query += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, *role)
		argIndex++
	}

	if isActive != nil {
		query += fmt.Sprintf(" AND is_active = $%d", argIndex)
		args = append(args, *isActive)
		argIndex++
	}

	query += " ORDER BY name"

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		var passwordHash sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.EmployeeID,
			&user.Name,
			&user.Email,
			&user.Department,
			&user.Role,
			&passwordHash,
			&user.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if passwordHash.Valid {
			user.PasswordHash = passwordHash.String
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %w", err)
	}

	return users, nil
}

// Update updates an existing user in the database
func (r *PGUserRepository) Update(user *entities.User) error {
	query := `
		UPDATE users
		SET employee_id = $2, name = $3, email = $4, department = $5, role = $6, is_active = $7
		WHERE id = $1
	`

	_, err := r.db.Exec(context.Background(), query,
		user.GetID(),
		user.GetEmployeeID(),
		user.GetName(),
		user.GetEmail(),
		user.GetDepartment(),
		user.GetRole(),
		user.IsActive(),
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete removes a user from the database (soft delete by deactivating)
func (r *PGUserRepository) Delete(id string) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Exists checks if a user exists with the given email or employee ID
func (r *PGUserRepository) Exists(email, employeeID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users
			WHERE (email = $1 OR employee_id = $2) AND is_active = true
		)
	`

	var exists bool
	err := r.db.QueryRow(context.Background(), query, email, employeeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// Create a struct to hold the scanned database user fields
type dbUser struct {
	ID           uuid.UUID      `db:"id"`
	EmployeeID   string         `db:"employee_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Department   sql.NullString `db:"department"`
	Role         string         `db:"role"`
	PasswordHash string         `db:"password_hash"`
	IsActive     bool           `db:"is_active"`
}

// ToDomainEntity converts a database user to a domain entity
func (du *dbUser) ToDomainEntity() (*entities.User, error) {
	passwordHasher := auth.NewPasswordHasher()

	department := ""
	if du.Department.Valid {
		department = du.Department.String
	}

	user, err := entities.NewUser(
		du.EmployeeID,
		du.Name,
		du.Email,
		department,
		entities.UserRole(du.Role),
		"", // No password needed when loading from DB
		passwordHasher,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Set the ID and other fields that are not set by NewUser
	user.SetID(du.ID.String())
	user.SetPasswordHash(du.PasswordHash)

	return user, nil
}