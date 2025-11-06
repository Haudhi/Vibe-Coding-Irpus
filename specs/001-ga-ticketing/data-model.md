# Domain Data Model: GA Ticketing System

**Branch**: `001-ga-ticketing` | **Date**: 2025-11-06 | **Spec**: [spec.md](./spec.md)

## Core Domain Entities

### 1. User (Aggregate Root)

```go
type User struct {
    ID           string    `json:"id" db:"id"`
    EmployeeID   string    `json:"employee_id" db:"employee_id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    Department   string    `json:"department" db:"department"`
    Role         UserRole  `json:"role" db:"role"`
    IsActive     bool      `json:"is_active" db:"is_active"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
    UserRoleRequester UserRole = "requester"  // Employee
    UserRoleApprover  UserRole = "approver"   // Budget approver
    UserRoleAdmin     UserRole = "admin"      // GA admin
)
```

### 2. Ticket (Aggregate Root)

```go
type Ticket struct {
    ID              string          `json:"id" db:"id"`
    TicketNumber    string          `json:"ticket_number" db:"ticket_number"`
    Title           string          `json:"title" db:"title"`
    Description     string          `json:"description" db:"description"`
    Category        TicketCategory  `json:"category" db:"category"`
    Priority        TicketPriority  `json:"priority" db:"priority"`
    Status          TicketStatus    `json:"status" db:"status"`
    RequesterID     string          `json:"requester_id" db:"requester_id"`
    AssignedAdminID *string         `json:"assigned_admin_id" db:"assigned_admin_id"`
    EstimatedCost   int64           `json:"estimated_cost" db:"estimated_cost"`   // in Rupiah
    ActualCost      *int64          `json:"actual_cost" db:"actual_cost"`         // in Rupiah
    RequiresApproval bool           `json:"requires_approval" db:"requires_approval"`
    CreatedAt       time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
    CompletedAt     *time.Time      `json:"completed_at" db:"completed_at"`

    // Associations
    Requester       User            `json:"requester" db:"-"`
    AssignedAdmin   *User           `json:"assigned_admin" db:"-"`
    Comments        []Comment       `json:"comments" db:"-"`
    StatusHistory   []StatusHistory `json:"status_history" db:"-"`
    Approvals       []Approval      `json:"approvals" db:"-"`
}

type TicketCategory string

const (
    CategoryOfficeSupplies     TicketCategory = "office_supplies"
    CategoryFacilityMaintenance TicketCategory = "facility_maintenance"
    CategoryPantrySupplies     TicketCategory = "pantry_supplies"
    CategoryMeetingRoom        TicketCategory = "meeting_room"
    CategoryOfficeFurniture    TicketCategory = "office_furniture"
    CategoryGeneralService     TicketCategory = "general_service"
)

type TicketPriority string

const (
    PriorityLow    TicketPriority = "low"
    PriorityMedium TicketPriority = "medium"
    PriorityHigh   TicketPriority = "high"
)

type TicketStatus string

const (
    StatusPending          TicketStatus = "pending"
    StatusWaitingApproval  TicketStatus = "waiting_approval"
    StatusApproved         TicketStatus = "approved"
    StatusRejected         TicketStatus = "rejected"
    StatusInProgress       TicketStatus = "in_progress"
    StatusCompleted        TicketStatus = "completed"
    StatusClosed           TicketStatus = "closed"
)
```

### 3. Asset (Aggregate Root)

```go
type Asset struct {
    ID                string      `json:"id" db:"id"`
    AssetCode         string      `json:"asset_code" db:"asset_code"`
    Name              string      `json:"name" db:"name"`
    Description       string      `json:"description" db:"description"`
    Category          AssetCategory `json:"category" db:"category"`
    Quantity          int         `json:"quantity" db:"quantity"`
    AvailableQuantity int         `json:"available_quantity" db:"available_quantity"`
    Location          string      `json:"location" db:"location"`
    Condition         AssetCondition `json:"condition" db:"condition"`
    UnitCost          int64       `json:"unit_cost" db:"unit_cost"`     // in Rupiah
    LastMaintenanceAt *time.Time  `json:"last_maintenance_at" db:"last_maintenance_at"`
    NextMaintenanceAt *time.Time  `json:"next_maintenance_at" db:"next_maintenance_at"`
    CreatedAt         time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`

    // Associations
    InventoryLogs     []InventoryLog `json:"inventory_logs" db:"-"`
}

type AssetCategory string

const (
    AssetCategoryOfficeFurniture    AssetCategory = "office_furniture"
    AssetCategoryOfficeSupplies     AssetCategory = "office_supplies"
    AssetCategoryPantrySupplies     AssetCategory = "pantry_supplies"
    AssetCategoryFacilityEquipment  AssetCategory = "facility_equipment"
    AssetCategoryMeetingRoom        AssetCategory = "meeting_room_equipment"
    AssetCategoryCleaningSupplies   AssetCategory = "cleaning_supplies"
)

type AssetCondition string

const (
    ConditionGood            AssetCondition = "good"
    ConditionNeedsMaintenance AssetCondition = "needs_maintenance"
    ConditionBroken          AssetCondition = "broken"
)
```

### 4. Comment (Entity)

```go
type Comment struct {
    ID        string    `json:"id" db:"id"`
    TicketID  string    `json:"ticket_id" db:"ticket_id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`

    // Associations
    User      User      `json:"user" db:"-"`
}
```

### 5. Approval (Entity)

```go
type Approval struct {
    ID          string       `json:"id" db:"id"`
    TicketID    string       `json:"ticket_id" db:"ticket_id"`
    ApproverID  string       `json:"approver_id" db:"approver_id"`
    Status      ApprovalStatus `json:"status" db:"status"`
    Comments    string       `json:"comments" db:"comments"`
    CreatedAt   time.Time    `json:"created_at" db:"created_at"`

    // Associations
    Approver    User         `json:"approver" db:"-"`
}

type ApprovalStatus string

const (
    ApprovalStatusPending  ApprovalStatus = "pending"
    ApprovalStatusApproved ApprovalStatus = "approved"
    ApprovalStatusRejected ApprovalStatus = "rejected"
)
```

### 6. StatusHistory (Entity)

```go
type StatusHistory struct {
    ID          string      `json:"id" db:"id"`
    TicketID    string      `json:"ticket_id" db:"ticket_id"`
    FromStatus  TicketStatus `json:"from_status" db:"from_status"`
    ToStatus    TicketStatus `json:"to_status" db:"to_status"`
    ChangedBy   string      `json:"changed_by" db:"changed_by"`
    Comments    string      `json:"comments" db:"comments"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`

    // Associations
    ChangedByUser User       `json:"changed_by_user" db:"-"`
}
```

### 7. InventoryLog (Entity)

```go
type InventoryLog struct {
    ID          string      `json:"id" db:"id"`
    AssetID     string      `json:"asset_id" db:"asset_id"`
    ChangeType  ChangeType  `json:"change_type" db:"change_type"`
    Quantity    int         `json:"quantity" db:"quantity"`
    Reason      string      `json:"reason" db:"reason"`
    CreatedBy   string      `json:"created_by" db:"created_by"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`

    // Associations
    CreatedByUser User       `json:"created_by_user" db:"-"`
}

type ChangeType string

const (
    ChangeTypeAdd    ChangeType = "add"
    ChangeTypeRemove ChangeType = "remove"
    ChangeTypeAdjust ChangeType = "adjust"
)
```

## Domain Services and Business Rules

### TicketService Business Rules

```go
// Approval rules validation
func (s *TicketService) RequiresApproval(category TicketCategory, estimatedCost int64) bool {
    return category == CategoryOfficeFurniture || estimatedCost >= 500000
}

// Ticket number generation
func (s *TicketService) GenerateTicketNumber() string {
    year := time.Now().Year()
    sequence := s.ticketRepository.GetNextSequence(year)
    return fmt.Sprintf("GA-%d-%04d", year, sequence)
}

// Status transition validation
func (s *TicketService) CanTransition(from, to TicketStatus) bool {
    validTransitions := map[TicketStatus][]TicketStatus{
        StatusPending:          {StatusWaitingApproval, StatusInProgress, StatusClosed},
        StatusWaitingApproval:  {StatusApproved, StatusRejected, StatusClosed},
        StatusApproved:         {StatusInProgress, StatusClosed},
        StatusRejected:         {StatusClosed},
        StatusInProgress:       {StatusCompleted, StatusClosed},
        StatusCompleted:        {StatusClosed},
        StatusClosed:           {},
    }

    allowed, exists := validTransitions[from]
    if !exists {
        return false
    }

    for _, status := range allowed {
        if status == to {
            return true
        }
    }
    return false
}
```

### AssetService Business Rules

```go
// Inventory availability check
func (s *AssetService) CheckAvailability(assetID string, quantity int) error {
    asset, err := s.assetRepository.FindByID(assetID)
    if err != nil {
        return err
    }

    if asset.AvailableQuantity < quantity {
        return errors.New("insufficient stock")
    }

    return nil
}

// Maintenance scheduling
func (s *AssetService) ScheduleMaintenance(assetID string, lastMaintenance time.Time) time.Time {
    // Schedule next maintenance 6 months from last maintenance
    return lastMaintenance.AddDate(0, 6, 0)
}
```

## Value Objects

### Money

```go
type Money struct {
    Amount   int64 `json:"amount"`
    Currency string `json:"currency"`
}

func NewMoney(amount int64) Money {
    return Money{
        Amount:   amount,
        Currency: "IDR",
    }
}

func (m Money) String() string {
    return fmt.Sprintf("Rp %d", m.Amount)
}
```

### Email

```go
type Email struct {
    Address string
}

func NewEmail(address string) (Email, error) {
    if !isValidEmail(address) {
        return Email{}, errors.New("invalid email format")
    }
    return Email{Address: address}, nil
}
```

## Aggregates and Bounded Contexts

### Ticket Management Bounded Context
- **Aggregate Root**: Ticket
- **Entities**: Comment, Approval, StatusHistory
- **Value Objects**: Money
- **Domain Services**: TicketService, ApprovalService

### Asset Management Bounded Context
- **Aggregate Root**: Asset
- **Entities**: InventoryLog
- **Value Objects**: AssetCode, Location
- **Domain Services**: AssetService, InventoryService

### User Management Bounded Context
- **Aggregate Root**: User
- **Domain Services**: AuthenticationService, AuthorizationService

## Database Schema Relationships

```
Users (1) ←→ (N) Tickets (requester_id)
Users (1) ←→ (N) Tickets (assigned_admin_id)
Tickets (1) ←→ (N) Comments
Tickets (1) ←→ (N) Approvals
Tickets (1) ←→ (N) StatusHistory

Assets (1) ←→ (N) InventoryLogs
Users (1) ←→ (N) InventoryLogs (created_by)

Users (1) ←→ (N) Comments (user_id)
Users (1) ←→ (N) Approvals (approver_id)
Users (1) ←→ (N) StatusHistory (changed_by)
```

## Indexes and Performance

### Recommended Database Indexes

```sql
-- Tickets
CREATE INDEX idx_tickets_requester_id ON tickets(requester_id);
CREATE INDEX idx_tickets_assigned_admin_id ON tickets(assigned_admin_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_category ON tickets(category);
CREATE INDEX idx_tickets_created_at ON tickets(created_at);
CREATE INDEX idx_tickets_ticket_number ON tickets(ticket_number) UNIQUE;

-- Assets
CREATE INDEX idx_assets_category ON assets(category);
CREATE INDEX idx_assets_location ON assets(location);
CREATE INDEX idx_assets_condition ON assets(condition);
CREATE INDEX idx_assets_asset_code ON assets(asset_code) UNIQUE;

-- Comments
CREATE INDEX idx_comments_ticket_id ON comments(ticket_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- Approvals
CREATE INDEX idx_approvals_ticket_id ON approvals(ticket_id);
CREATE INDEX idx_approvals_approver_id ON approvals(approver_id);

-- Status History
CREATE INDEX idx_status_history_ticket_id ON status_history(ticket_id);
CREATE INDEX idx_status_history_created_at ON status_history(created_at);

-- Inventory Logs
CREATE INDEX idx_inventory_logs_asset_id ON inventory_logs(asset_id);
CREATE INDEX idx_inventory_logs_created_at ON inventory_logs(created_at);
```

## Data Validation Constraints

```go
type TicketCreateRequest struct {
    Title          string          `json:"title" validate:"required,max=255"`
    Description    string          `json:"description" validate:"required"`
    Category       TicketCategory  `json:"category" validate:"required,oneof=office_supplies facility_maintenance pantry_supplies meeting_room office_furniture general_service"`
    Priority       TicketPriority  `json:"priority" validate:"required,oneof=low medium high"`
    EstimatedCost  int64           `json:"estimated_cost" validate:"required,min=0"`
}

type AssetCreateRequest struct {
    Name         string        `json:"name" validate:"required,max=255"`
    Description  string        `json:"description" validate:"required"`
    Category     AssetCategory `json:"category" validate:"required,oneof=office_furniture office_supplies pantry_supplies facility_equipment meeting_room_equipment cleaning_supplies"`
    Quantity     int           `json:"quantity" validate:"required,min=1"`
    Location     string        `json:"location" validate:"required,max=255"`
    UnitCost     int64         `json:"unit_cost" validate:"required,min=0"`
}
```

## Authentication and Authorization

### JWT Token Structure

```go
type Claims struct {
    UserID      string    `json:"user_id"`
    EmployeeID  string    `json:"employee_id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    Role        UserRole  `json:"role"`
    Department  string    `json:"department"`
    jwt.RegisteredClaims
}
```

### Permissions Matrix

| Operation | Requester | Approver | Admin |
|-----------|-----------|----------|-------|
| Create Ticket | ✅ | ✅ | ✅ |
| View Own Tickets | ✅ | ✅ | ✅ |
| View All Tickets | ❌ | ❌ | ✅ |
| Assign Ticket | ❌ | ❌ | ✅ |
| Approve Ticket | ❌ | ✅ | ✅ |
| Reject Ticket | ❌ | ✅ | ✅ |
| Manage Assets | ❌ | ❌ | ✅ |
| Update Inventory | ❌ | ❌ | ✅ |
| Add Comments | ✅* | ✅* | ✅* |

*Only on tickets they have access to (own, assigned, or approval-eligible)