package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/repositories"
)

// PGApprovalRepository implements the ApprovalRepository interface for PostgreSQL
type PGApprovalRepository struct {
	db *pgxpool.Pool
}

// NewPGApprovalRepository creates a new PostgreSQL approval repository
func NewPGApprovalRepository(db *pgxpool.Pool) repositories.ApprovalRepository {
	return &PGApprovalRepository{db: db}
}

// Create saves a new approval to the database
func (r *PGApprovalRepository) Create(approval *entities.Approval) error {
	query := `
		INSERT INTO approvals (id, ticket_id, approver_id, status, comments)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(context.Background(), query,
		approval.GetID(),
		approval.GetTicketID(),
		approval.GetApproverID(),
		approval.GetStatus(),
		approval.GetComments(),
	)

	if err != nil {
		return fmt.Errorf("failed to create approval: %w", err)
	}

	return nil
}

// FindByID retrieves an approval by their ID
func (r *PGApprovalRepository) FindByID(id string) (*entities.Approval, error) {
	query := `
		SELECT id, ticket_id, approver_id, status, comments, created_at
		FROM approvals
		WHERE id = $1
	`

	var approval entities.Approval
	var comments sql.NullString

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&approval.ID,
		&approval.TicketID,
		&approval.ApproverID,
		&approval.Status,
		&comments,
		&approval.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("approval not found")
		}
		return nil, fmt.Errorf("failed to find approval: %w", err)
	}

	if comments.Valid {
		approval.Comments = comments.String
	}

	return &approval, nil
}

// FindByTicketID retrieves all approvals for a specific ticket
func (r *PGApprovalRepository) FindByTicketID(ticketID string) ([]*entities.Approval, error) {
	query := `
		SELECT id, ticket_id, approver_id, status, comments, created_at
		FROM approvals
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approvals: %w", err)
	}
	defer rows.Close()

	var approvals []*entities.Approval
	for rows.Next() {
		var approval entities.Approval
		var comments sql.NullString

		err := rows.Scan(
			&approval.ID,
			&approval.TicketID,
			&approval.ApproverID,
			&approval.Status,
			&comments,
			&approval.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}

		if comments.Valid {
			approval.Comments = comments.String
		}

		approvals = append(approvals, &approval)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating approvals: %w", err)
	}

	return approvals, nil
}

// FindByApproverID retrieves all approvals by a specific approver
func (r *PGApprovalRepository) FindByApproverID(approverID string) ([]*entities.Approval, error) {
	query := `
		SELECT id, ticket_id, approver_id, status, comments, created_at
		FROM approvals
		WHERE approver_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query approvals: %w", err)
	}
	defer rows.Close()

	var approvals []*entities.Approval
	for rows.Next() {
		var approval entities.Approval
		var comments sql.NullString

		err := rows.Scan(
			&approval.ID,
			&approval.TicketID,
			&approval.ApproverID,
			&approval.Status,
			&comments,
			&approval.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}

		if comments.Valid {
			approval.Comments = comments.String
		}

		approvals = append(approvals, &approval)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating approvals: %w", err)
	}

	return approvals, nil
}

// FindPendingByApprover retrieves all pending approvals for a specific approver
func (r *PGApprovalRepository) FindPendingByApprover(approverID string) ([]*entities.Approval, error) {
	query := `
		SELECT id, ticket_id, approver_id, status, comments, created_at
		FROM approvals
		WHERE approver_id = $1 AND status = 'pending'
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, approverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending approvals: %w", err)
	}
	defer rows.Close()

	var approvals []*entities.Approval
	for rows.Next() {
		var approval entities.Approval
		var comments sql.NullString

		err := rows.Scan(
			&approval.ID,
			&approval.TicketID,
			&approval.ApproverID,
			&approval.Status,
			&comments,
			&approval.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}

		if comments.Valid {
			approval.Comments = comments.String
		}

		approvals = append(approvals, &approval)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating approvals: %w", err)
	}

	return approvals, nil
}

// FindPending retrieves all pending approvals in the system
func (r *PGApprovalRepository) FindPending() ([]*entities.Approval, error) {
	query := `
		SELECT id, ticket_id, approver_id, status, comments, created_at
		FROM approvals
		WHERE status = 'pending'
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pending approvals: %w", err)
	}
	defer rows.Close()

	var approvals []*entities.Approval
	for rows.Next() {
		var approval entities.Approval
		var comments sql.NullString

		err := rows.Scan(
			&approval.ID,
			&approval.TicketID,
			&approval.ApproverID,
			&approval.Status,
			&comments,
			&approval.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan approval: %w", err)
		}

		if comments.Valid {
			approval.Comments = comments.String
		}

		approvals = append(approvals, &approval)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating approvals: %w", err)
	}

	return approvals, nil
}

// FindByTicketAndApprover retrieves an approval for a specific ticket and approver
func (r *PGApprovalRepository) FindByTicketAndApprover(ticketID, approverID string) (*entities.Approval, error) {
	query := `
		SELECT id, ticket_id, approver_id, status, comments, created_at
		FROM approvals
		WHERE ticket_id = $1 AND approver_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`

	var approval entities.Approval
	var comments sql.NullString

	err := r.db.QueryRow(context.Background(), query, ticketID, approverID).Scan(
		&approval.ID,
		&approval.TicketID,
		&approval.ApproverID,
		&approval.Status,
		&comments,
		&approval.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("approval not found")
		}
		return nil, fmt.Errorf("failed to find approval: %w", err)
	}

	if comments.Valid {
		approval.Comments = comments.String
	}

	return &approval, nil
}

// Update updates an existing approval in the database
func (r *PGApprovalRepository) Update(approval *entities.Approval) error {
	query := `
		UPDATE approvals
		SET status = $2, comments = $3
		WHERE id = $1
	`

	_, err := r.db.Exec(context.Background(), query,
		approval.GetID(),
		approval.GetStatus(),
		approval.GetComments(),
	)

	if err != nil {
		return fmt.Errorf("failed to update approval: %w", err)
	}

	return nil
}

// Delete removes an approval from the database
func (r *PGApprovalRepository) Delete(id string) error {
	query := `DELETE FROM approvals WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete approval: %w", err)
	}

	return nil
}

// CheckExists checks if an approval exists for a specific ticket and approver
func (r *PGApprovalRepository) CheckExists(ticketID, approverID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM approvals
			WHERE ticket_id = $1 AND approver_id = $2
		)
	`

	var exists bool
	err := r.db.QueryRow(context.Background(), query, ticketID, approverID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check approval existence: %w", err)
	}

	return exists, nil
}