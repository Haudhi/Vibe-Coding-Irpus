package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/repositories"
)

// PGCommentRepository implements the CommentRepository interface for PostgreSQL
type PGCommentRepository struct {
	db *pgxpool.Pool
}

// NewPGCommentRepository creates a new PostgreSQL comment repository
func NewPGCommentRepository(db *pgxpool.Pool) repositories.CommentRepository {
	return &PGCommentRepository{db: db}
}

// Create saves a new comment to the database
func (r *PGCommentRepository) Create(comment *entities.Comment) error {
	query := `
		INSERT INTO comments (id, ticket_id, user_id, content)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(context.Background(), query,
		comment.GetID(),
		comment.GetTicketID(),
		comment.GetUserID(),
		comment.GetContent(),
	)

	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// FindByID retrieves a comment by their ID
func (r *PGCommentRepository) FindByID(id string) (*entities.Comment, error) {
	query := `
		SELECT id, ticket_id, user_id, content, created_at
		FROM comments
		WHERE id = $1
	`

	var comment entities.Comment

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&comment.ID,
		&comment.TicketID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	return &comment, nil
}

// FindByTicketID retrieves all comments for a specific ticket
func (r *PGCommentRepository) FindByTicketID(ticketID string) ([]*entities.Comment, error) {
	query := `
		SELECT id, ticket_id, user_id, content, created_at
		FROM comments
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(context.Background(), query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		var comment entities.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.TicketID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

// FindByUserID retrieves all comments by a specific user
func (r *PGCommentRepository) FindByUserID(userID string) ([]*entities.Comment, error) {
	query := `
		SELECT id, ticket_id, user_id, content, created_at
		FROM comments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		var comment entities.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.TicketID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}

		comments = append(comments, &comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %w", err)
	}

	return comments, nil
}

// Update updates an existing comment in the database
func (r *PGCommentRepository) Update(comment *entities.Comment) error {
	query := `
		UPDATE comments
		SET content = $2
		WHERE id = $1
	`

	_, err := r.db.Exec(context.Background(), query,
		comment.GetID(),
		comment.GetContent(),
	)

	if err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}

	return nil
}

// Delete removes a comment from the database
func (r *PGCommentRepository) Delete(id string) error {
	query := `DELETE FROM comments WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

// DeleteByTicketID removes all comments for a specific ticket
func (r *PGCommentRepository) DeleteByTicketID(ticketID string) error {
	query := `DELETE FROM comments WHERE ticket_id = $1`

	_, err := r.db.Exec(context.Background(), query, ticketID)
	if err != nil {
		return fmt.Errorf("failed to delete comments for ticket: %w", err)
	}

	return nil
}