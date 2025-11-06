package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// ticketRecord represents the database model for a ticket
type ticketRecord struct {
	ID               string
	TicketNumber     string
	Title            string
	Description      string
	Category         string
	Priority         string
	Status           string
	RequesterID      string
	AssignedAdminID  sql.NullString
	EstimatedCost    sql.NullInt64
	ActualCost       sql.NullInt64
	RequiresApproval bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CompletedAt      sql.NullTime
	AssignedAt       sql.NullTime
}

// TicketRepository implements the TicketRepository interface using PostgreSQL
type TicketRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewTicketRepository creates a new TicketRepository
func NewTicketRepository(pool *pgxpool.Pool, logger *zap.Logger) services.TicketRepository {
	return &TicketRepository{
		pool:   pool,
		logger: logger,
	}
}

// Create saves a new ticket to the database
func (r *TicketRepository) Create(ctx context.Context, ticket *entities.Ticket) error {
	query := `
		INSERT INTO tickets (
			id, ticket_number, title, description, category, priority,
			status, requester_id, assigned_admin_id, estimated_cost,
			actual_cost, requires_approval, created_at, updated_at, completed_at, assigned_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	var completedAt, assignedAt sql.NullTime
	if ticket.GetCompletedAt() != nil {
		completedAt = sql.NullTime{Time: *ticket.GetCompletedAt(), Valid: true}
	}
	if ticket.GetAssignedAt() != nil {
		assignedAt = sql.NullTime{Time: *ticket.GetAssignedAt(), Valid: true}
	}

	_, err := r.pool.Exec(ctx, query,
		ticket.GetID(),
		ticket.GetTicketNumber(),
		ticket.GetTitle(),
		ticket.GetDescription(),
		string(ticket.GetCategory()),
		string(ticket.GetPriority()),
		string(ticket.GetStatus()),
		ticket.GetRequesterID(),
		ticket.GetAssignedAdminID(),
		ticket.GetEstimatedCost().Amount,
		nil, // actual_cost will be set later
		ticket.RequiresApproval(),
		ticket.GetCreatedAt(),
		ticket.GetUpdatedAt(),
		completedAt,
		assignedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create ticket",
			zap.String("ticket_id", ticket.GetID()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	// Insert status history
	if err := r.insertStatusHistory(ctx, ticket); err != nil {
		r.logger.Error("Failed to insert status history",
			zap.String("ticket_id", ticket.GetID()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to insert status history: %w", err)
	}

	r.logger.Info("Ticket created successfully",
		zap.String("ticket_id", ticket.GetID()),
		zap.String("ticket_number", ticket.GetTicketNumber()),
	)

	return nil
}

// GetByID retrieves a ticket by ID
func (r *TicketRepository) GetByID(ctx context.Context, id string) (*entities.Ticket, error) {
	query := `
		SELECT
			id, ticket_number, title, description, category, priority,
			status, requester_id, assigned_admin_id, estimated_cost,
			actual_cost, requires_approval, created_at, updated_at, completed_at, assigned_at
		FROM tickets
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var ticket ticketRecord

	err := row.Scan(
		&ticket.ID,
		&ticket.TicketNumber,
		&ticket.Title,
		&ticket.Description,
		&ticket.Category,
		&ticket.Priority,
		&ticket.Status,
		&ticket.RequesterID,
		&ticket.AssignedAdminID,
		&ticket.EstimatedCost,
		&ticket.ActualCost,
		&ticket.RequiresApproval,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.CompletedAt,
		&ticket.AssignedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		r.logger.Error("Failed to get ticket by ID",
			zap.String("ticket_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Convert to domain entity
	domainTicket := r.mapToDomainEntity(&ticket)

	// TODO: Load status history and comments
	// These methods need to be implemented on the Ticket entity
	// statusHistory, err := r.getStatusHistory(ctx, id)
	// if err != nil {
	// 	r.logger.Warn("Failed to load status history",
	// 		zap.String("ticket_id", id),
	// 		zap.Error(err),
	// 	)
	// }
	// domainTicket.SetStatusHistory(statusHistory)

	// comments, err := r.getComments(ctx, id)
	// if err != nil {
	// 	r.logger.Warn("Failed to load comments",
	// 		zap.String("ticket_id", id),
	// 		zap.Error(err),
	// 	)
	// }
	// domainTicket.SetComments(comments)

	return domainTicket, nil
}

// GetByTicketNumber retrieves a ticket by ticket number
func (r *TicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*entities.Ticket, error) {
	query := `
		SELECT
			id, ticket_number, title, description, category, priority,
			status, requester_id, assigned_admin_id, estimated_cost,
			actual_cost, requires_approval, created_at, updated_at, completed_at, assigned_at
		FROM tickets
		WHERE ticket_number = $1
	`

	row := r.pool.QueryRow(ctx, query, ticketNumber)

	var ticket ticketRecord

	err := row.Scan(
		&ticket.ID,
		&ticket.TicketNumber,
		&ticket.Title,
		&ticket.Description,
		&ticket.Category,
		&ticket.Priority,
		&ticket.Status,
		&ticket.RequesterID,
		&ticket.AssignedAdminID,
		&ticket.EstimatedCost,
		&ticket.ActualCost,
		&ticket.RequiresApproval,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.CompletedAt,
		&ticket.AssignedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("ticket not found")
		}
		r.logger.Error("Failed to get ticket by ticket number",
			zap.String("ticket_number", ticketNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Convert to domain entity
	domainTicket := r.mapToDomainEntity(&ticket)

	return domainTicket, nil
}

// GetByRequesterID retrieves tickets by requester ID with pagination
func (r *TicketRepository) GetByRequesterID(ctx context.Context, requesterID string, limit, offset int) ([]*entities.Ticket, error) {
	query := `
		SELECT
			id, ticket_number, title, description, category, priority,
			status, requester_id, assigned_admin_id, estimated_cost,
			actual_cost, requires_approval, created_at, updated_at, completed_at, assigned_at
		FROM tickets
		WHERE requester_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, requesterID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get tickets by requester ID",
			zap.String("requester_id", requesterID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer rows.Close()

	var tickets []*entities.Ticket
	for rows.Next() {
		var ticket ticketRecord

		if err := rows.Scan(
			&ticket.ID,
			&ticket.TicketNumber,
			&ticket.Title,
			&ticket.Description,
			&ticket.Category,
			&ticket.Priority,
			&ticket.Status,
			&ticket.RequesterID,
			&ticket.AssignedAdminID,
			&ticket.EstimatedCost,
			&ticket.ActualCost,
			&ticket.RequiresApproval,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.CompletedAt,
			&ticket.AssignedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ticket row: %w", err)
		}

		domainTicket := r.mapToDomainEntity(&ticket)
		tickets = append(tickets, domainTicket)
	}

	return tickets, nil
}

// GetAll retrieves all tickets with pagination
func (r *TicketRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Ticket, error) {
	query := `
		SELECT
			id, ticket_number, title, description, category, priority,
			status, requester_id, assigned_admin_id, estimated_cost,
			actual_cost, requires_approval, created_at, updated_at, completed_at, assigned_at
		FROM tickets
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get all tickets", zap.Error(err))
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer rows.Close()

	var tickets []*entities.Ticket
	for rows.Next() {
		var ticket ticketRecord

		if err := rows.Scan(
			&ticket.ID,
			&ticket.TicketNumber,
			&ticket.Title,
			&ticket.Description,
			&ticket.Category,
			&ticket.Priority,
			&ticket.Status,
			&ticket.RequesterID,
			&ticket.AssignedAdminID,
			&ticket.EstimatedCost,
			&ticket.ActualCost,
			&ticket.RequiresApproval,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.CompletedAt,
			&ticket.AssignedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ticket row: %w", err)
		}

		domainTicket := r.mapToDomainEntity(&ticket)
		tickets = append(tickets, domainTicket)
	}

	return tickets, nil
}

// Update updates an existing ticket
func (r *TicketRepository) Update(ctx context.Context, ticket *entities.Ticket) error {
	query := `
		UPDATE tickets SET
			title = $2,
			description = $3,
			priority = $4,
			status = $5,
			assigned_admin_id = $6,
			actual_cost = $7,
			updated_at = $8,
			completed_at = $9,
			assigned_at = $10
		WHERE id = $1
	`

	var actualCost *int64
	if ticket.GetActualCost() != nil {
		actualCost = &ticket.GetActualCost().Amount
	}

	var completedAt, assignedAt sql.NullTime
	if ticket.GetCompletedAt() != nil {
		completedAt = sql.NullTime{Time: *ticket.GetCompletedAt(), Valid: true}
	}
	if ticket.GetAssignedAt() != nil {
		assignedAt = sql.NullTime{Time: *ticket.GetAssignedAt(), Valid: true}
	}

	_, err := r.pool.Exec(ctx, query,
		ticket.GetID(),
		ticket.GetTitle(),
		ticket.GetDescription(),
		string(ticket.GetPriority()),
		string(ticket.GetStatus()),
		ticket.GetAssignedAdminID(),
		actualCost,
		ticket.GetUpdatedAt(),
		completedAt,
		assignedAt,
	)

	if err != nil {
		r.logger.Error("Failed to update ticket",
			zap.String("ticket_id", ticket.GetID()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	// Update status history if needed
	if err := r.updateStatusHistory(ctx, ticket); err != nil {
		r.logger.Warn("Failed to update status history",
			zap.String("ticket_id", ticket.GetID()),
			zap.Error(err),
		)
	}

	r.logger.Info("Ticket updated successfully",
		zap.String("ticket_id", ticket.GetID()),
		zap.String("status", string(ticket.GetStatus())),
	)

	return nil
}

// Delete soft deletes a ticket by setting status to closed
func (r *TicketRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE tickets SET status = 'closed', updated_at = $1 WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		r.logger.Error("Failed to delete ticket",
			zap.String("ticket_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	r.logger.Info("Ticket deleted successfully",
		zap.String("ticket_id", id),
	)

	return nil
}

// GetNextSequenceNumber gets the next sequence number for ticket numbering
func (r *TicketRepository) GetNextSequenceNumber(ctx context.Context, year int) (int, error) {
	query := `
		SELECT COALESCE(MAX(CAST(SUBSTRING(ticket_number, 7, 4) AS INTEGER)), 0) + 1
		FROM tickets
		WHERE ticket_number LIKE $1
	`

	pattern := fmt.Sprintf("GA-%d-%%", year)
	row := r.pool.QueryRow(ctx, query, pattern)

	var sequence int
	if err := row.Scan(&sequence); err != nil {
		if err == pgx.ErrNoRows {
			return 1, nil
		}
		return 0, fmt.Errorf("failed to get sequence number: %w", err)
	}

	return sequence, nil
}

// Helper methods

func (r *TicketRepository) mapToDomainEntity(record *ticketRecord) *entities.Ticket {
	// TODO: This is a simplified mapping - in a real implementation,
	// you would need to properly reconstruct the domain entity from the database record
	// For now, return nil as placeholder
	_ = record
	return nil // nolint
}

func (r *TicketRepository) insertStatusHistory(ctx context.Context, ticket *entities.Ticket) error {
	query := `
		INSERT INTO status_history (id, ticket_id, from_status, to_status, changed_by, comments, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, history := range ticket.GetStatusHistory() {
		_, err := r.pool.Exec(ctx, query,
			history.ID,
			history.TicketID,
			history.FromStatus,
			history.ToStatus,
			history.ChangedBy,
			history.Comments,
			history.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert status history: %w", err)
		}
	}

	return nil
}

func (r *TicketRepository) updateStatusHistory(ctx context.Context, ticket *entities.Ticket) error {
	// Implementation for updating status history
	return nil
}

func (r *TicketRepository) getStatusHistory(ctx context.Context, ticketID string) ([]*entities.StatusHistory, error) {
	query := `
		SELECT id, ticket_id, from_status, to_status, changed_by, comments, created_at
		FROM status_history
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get status history: %w", err)
	}
	defer rows.Close()

	var history []*entities.StatusHistory
	for rows.Next() {
		var h entities.StatusHistory
		if err := rows.Scan(
			&h.ID,
			&h.TicketID,
			&h.FromStatus,
			&h.ToStatus,
			&h.ChangedBy,
			&h.Comments,
			&h.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan status history row: %w", err)
		}
		history = append(history, &h)
	}

	return history, nil
}

func (r *TicketRepository) getComments(ctx context.Context, ticketID string) ([]*entities.Comment, error) {
	// TODO: Implement comment loading - Comment entity has unexported fields
	// Need to either add getters/setters or use a constructor function
	_ = ticketID
	return nil, nil
}