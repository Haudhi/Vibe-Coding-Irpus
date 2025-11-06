package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApprovalStatus represents the status of an approval
type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusRejected ApprovalStatus = "rejected"
)

// Approval represents an approval decision for a ticket
type Approval struct {
	id         string
	ticketID   string
	approverID string
	status     ApprovalStatus
	comments   string
	createdAt  time.Time
}

// NewApproval creates a new approval
func NewApproval(ticketID, approverID string) (*Approval, error) {
	if ticketID == "" {
		return nil, errors.New("ticket ID is required")
	}
	if approverID == "" {
		return nil, errors.New("approver ID is required")
	}

	return &Approval{
		id:         uuid.New().String(),
		ticketID:   ticketID,
		approverID: approverID,
		status:     ApprovalStatusPending,
		createdAt:  time.Now(),
	}, nil
}

// Getters
func (a *Approval) GetID() string          { return a.id }
func (a *Approval) GetTicketID() string    { return a.ticketID }
func (a *Approval) GetApproverID() string  { return a.approverID }
func (a *Approval) GetStatus() ApprovalStatus { return a.status }
func (a *Approval) GetComments() string    { return a.comments }
func (a *Approval) GetCreatedAt() time.Time { return a.createdAt }

// Setters (business logic)
func (a *Approval) Approve(comments string) error {
	if a.status != ApprovalStatusPending {
		return fmt.Errorf("cannot approve approval with status: %s", a.status)
	}

	a.status = ApprovalStatusApproved
	a.comments = comments
	return nil
}

func (a *Approval) Reject(comments string) error {
	if a.status != ApprovalStatusPending {
		return fmt.Errorf("cannot reject approval with status: %s", a.status)
	}

	if comments == "" {
		return errors.New("rejection reason is required")
	}

	a.status = ApprovalStatusRejected
	a.comments = comments
	return nil
}

// IsPending checks if the approval is still pending
func (a *Approval) IsPending() bool {
	return a.status == ApprovalStatusPending
}

// IsApproved checks if the approval was approved
func (a *Approval) IsApproved() bool {
	return a.status == ApprovalStatusApproved
}

// IsRejected checks if the approval was rejected
func (a *Approval) IsRejected() bool {
	return a.status == ApprovalStatusRejected
}

// CanBeUpdated checks if the approval can still be updated
func (a *Approval) CanBeUpdated() bool {
	return a.status == ApprovalStatusPending
}

// ValidateStatus validates if the status is valid
func ValidateStatus(status string) (ApprovalStatus, error) {
	switch status {
	case string(ApprovalStatusPending):
		return ApprovalStatusPending, nil
	case string(ApprovalStatusApproved):
		return ApprovalStatusApproved, nil
	case string(ApprovalStatusRejected):
		return ApprovalStatusRejected, nil
	default:
		return "", fmt.Errorf("invalid approval status: %s", status)
	}
}

// GetAllApprovalStatuses returns all available approval statuses
func GetAllApprovalStatuses() []ApprovalStatus {
	return []ApprovalStatus{
		ApprovalStatusPending,
		ApprovalStatusApproved,
		ApprovalStatusRejected,
	}
}

// StatusDisplayName returns a human-readable display name for a status
func StatusDisplayName(status ApprovalStatus) string {
	switch status {
	case ApprovalStatusPending:
		return "Pending"
	case ApprovalStatusApproved:
		return "Approved"
	case ApprovalStatusRejected:
		return "Rejected"
	default:
		return string(status)
	}
}