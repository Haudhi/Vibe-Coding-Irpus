# GA Ticketing API Documentation - Competition Version (GA Focus Only)

## Base URL

```
http://localhost:8080/api
```

## User Roles

### 3 Roles in System:

1. **Requester** (Employee) - Can create and view own GA tickets
2. **Approver** - Can approve/reject GA requests that need approval
3. **Admin** - Full access to tickets, assets, inventory, and system management

---

## 2. GA Ticket APIs

### 2.1 Create GA Ticket

**Endpoint:** `POST /api/tickets`

**Description:** Create a new GA service request

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**

```json
{
  "title": "Request office stationery supplies",
  "description": "Need 5 boxes of A4 paper, 10 pens, and 3 staplers for Finance department",
  "priority": "medium",
  "category": "office_supplies",
  "estimated_cost": 250000
}
```

**GA Service Categories:**

- `office_supplies` - ATK (Alat Tulis Kantor) - pens, paper, staplers, etc
- `facility_maintenance` - AC repair, cleaning, building maintenance
- `pantry_supplies` - Coffee, tea, snacks, drinking water
- `meeting_room` - Meeting room booking and setup
- `office_furniture` - Chairs, desks, cabinets requests
- `general_service` - Other GA services

**Field Validations:**

- `title`: required, max 255 characters
- `description`: required
- `priority`: required, enum: `low`, `medium`, `high`
- `category`: required (see categories above)
- `estimated_cost`: optional, decimal (Rupiah)

**Success Response (201 Created):**

```json
{
  "id": 15,
  "ticket_number": "GA-2025-0015",
  "title": "Request office stationery supplies",
  "description": "Need 5 boxes of A4 paper, 10 pens, and 3 staplers for Finance department",
  "status": "pending",
  "priority": "medium",
  "category": "office_supplies",
  "estimated_cost": 250000,
  "requires_approval": false,
  "requester_id": 1,
  "requester_name": "John Doe",
  "requester_department": "Finance",
  "assigned_to_id": null,
  "assigned_to_name": null,
  "created_at": "2025-01-15T14:30:00Z",
  "updated_at": "2025-01-15T14:30:00Z"
}
```

**Approval Rules:**

- Requests with `estimated_cost` < Rp 500,000: No approval needed
- Requests with `estimated_cost` >= Rp 500,000: Requires 1 approval
- All `office_furniture` requests: Always require approval

**Error Response (400 Bad Request):**

```json
{
  "error": "title is required"
}
```

---

### 2.2 List GA Tickets

**Endpoint:** `GET /api/tickets`

**Description:** Get list of GA tickets (filtered by role)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Role-based Access:**

- **Requester**: Only sees own tickets
- **Approver**: Sees tickets needing approval
- **Admin**: Sees all tickets

**Query Parameters (Optional):**

```
?status=pending            # Filter by status
?priority=high             # Filter by priority
?category=office_supplies  # Filter by category
?requires_approval=true    # Filter tickets needing approval
```

**Example Request:**

```
GET /api/tickets?status=pending&category=office_supplies
```

**Success Response (200 OK):**

```json
{
  "tickets": [
    {
      "id": 1,
      "ticket_number": "GA-2025-0001",
      "title": "Request office stationery supplies",
      "description": "Need office supplies for Finance department",
      "status": "pending",
      "priority": "medium",
      "category": "office_supplies",
      "estimated_cost": 250000,
      "requires_approval": false,
      "requester_id": 1,
      "requester_name": "John Doe",
      "requester_department": "Finance",
      "assigned_to_id": 3,
      "assigned_to_name": "Admin GA",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "ticket_number": "GA-2025-0002",
      "title": "AC maintenance in meeting room",
      "description": "AC not cooling properly in Meeting Room A",
      "status": "in_progress",
      "priority": "high",
      "category": "facility_maintenance",
      "estimated_cost": 1500000,
      "requires_approval": true,
      "approval_status": "approved",
      "requester_id": 1,
      "requester_name": "John Doe",
      "requester_department": "Finance",
      "assigned_to_id": 3,
      "assigned_to_name": "Admin GA",
      "created_at": "2025-01-14T11:00:00Z",
      "updated_at": "2025-01-15T09:30:00Z"
    }
  ],
  "total": 2
}
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "Unauthorized"
}
```

---

### 2.3 Get GA Ticket Detail

**Endpoint:** `GET /api/tickets/:id`

**Description:** Get detailed information of a specific GA ticket

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Example Request:**

```
GET /api/tickets/1
```

**Success Response (200 OK):**

```json
{
  "id": 1,
  "ticket_number": "GA-2025-0001",
  "title": "Request office stationery supplies",
  "description": "Need 5 boxes of A4 paper, 10 pens, and 3 staplers for Finance department",
  "status": "completed",
  "priority": "medium",
  "category": "office_supplies",
  "estimated_cost": 250000,
  "actual_cost": 245000,
  "requires_approval": false,
  "requester": {
    "id": 1,
    "name": "John Doe",
    "email": "requester@company.com",
    "department": "Finance"
  },
  "assigned_to": {
    "id": 3,
    "name": "Admin GA",
    "email": "admin@company.com"
  },
  "approval_info": null,
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T16:30:00Z",
  "completed_at": "2025-01-15T16:30:00Z"
}
```

**With Approval Info (for tickets requiring approval):**

```json
{
  "id": 5,
  "ticket_number": "GA-2025-0005",
  "title": "Purchase 10 new office chairs",
  "description": "Need ergonomic office chairs for new employees",
  "status": "waiting_approval",
  "priority": "high",
  "category": "office_furniture",
  "estimated_cost": 15000000,
  "requires_approval": true,
  "requester": {
    "id": 1,
    "name": "John Doe",
    "email": "requester@company.com",
    "department": "Finance"
  },
  "assigned_to": null,
  "approval_info": {
    "status": "pending",
    "approver_id": 2,
    "approver_name": "Jane Approver",
    "approved_at": null,
    "approval_notes": null
  },
  "created_at": "2025-01-15T14:00:00Z",
  "updated_at": "2025-01-15T14:00:00Z"
}
```

**Error Response (404 Not Found):**

```json
{
  "error": "Ticket not found"
}
```

---

### 2.4 Update GA Ticket Status

**Endpoint:** `PUT /api/tickets/:id/status`

**Description:** Update ticket status (Admin only can assign and process)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**

```json
{
  "status": "in_progress",
  "assigned_to_id": 3,
  "actual_cost": 245000,
  "notes": "Supplies ordered from vendor"
}
```

**Ticket Status Workflow:**

- `pending` → `waiting_approval` (if requires approval)
- `pending` / `approved` → `in_progress` (Admin assigns to self)
- `in_progress` → `completed` (Admin completes)
- `completed` → `closed` (Requester confirms)

**Status Options:**

- `pending` - Ticket created, waiting for processing
- `waiting_approval` - Waiting for approver action
- `approved` - Approved by approver, ready to process
- `rejected` - Rejected by approver
- `in_progress` - Admin is working on it
- `completed` - Work completed
- `closed` - Requester confirmed completion

**Field Validations:**

- `status`: required
- `assigned_to_id`: optional (admin user ID)
- `actual_cost`: optional (final cost in Rupiah)
- `notes`: optional

**Success Response (200 OK):**

```json
{
  "id": 1,
  "ticket_number": "GA-2025-0001",
  "title": "Request office stationery supplies",
  "status": "in_progress",
  "assigned_to_id": 3,
  "assigned_to_name": "Admin GA",
  "updated_at": "2025-01-15T14:45:00Z",
  "message": "Ticket status updated successfully"
}
```

**Error Response (403 Forbidden):**

```json
{
  "error": "Only admin can update ticket status"
}
```

---

### 2.5 Approve/Reject GA Ticket

**Endpoint:** `PUT /api/tickets/:id/approval`

**Description:** Approve or reject GA ticket (Approver only)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**

```json
{
  "action": "approve",
  "notes": "Approved. Please proceed with purchase."
}
```

**Field Validations:**

- `action`: required, enum: `approve`, `reject`
- `notes`: optional, approval/rejection notes

**Example Request (Approve):**

```
PUT /api/tickets/5/approval
```

```json
{
  "action": "approve",
  "notes": "Budget approved for office chairs purchase"
}
```

**Success Response (200 OK):**

```json
{
  "id": 5,
  "ticket_number": "GA-2025-0005",
  "title": "Purchase 10 new office chairs",
  "status": "approved",
  "approval_info": {
    "status": "approved",
    "approver_id": 2,
    "approver_name": "Jane Approver",
    "approved_at": "2025-01-15T15:00:00Z",
    "approval_notes": "Budget approved for office chairs purchase"
  },
  "message": "Ticket approved successfully"
}
```

**Example Request (Reject):**

```json
{
  "action": "reject",
  "notes": "Budget not available this month. Please resubmit next month."
}
```

**Success Response (200 OK):**

```json
{
  "id": 5,
  "ticket_number": "GA-2025-0005",
  "title": "Purchase 10 new office chairs",
  "status": "rejected",
  "approval_info": {
    "status": "rejected",
    "approver_id": 2,
    "approver_name": "Jane Approver",
    "approved_at": "2025-01-15T15:00:00Z",
    "approval_notes": "Budget not available this month. Please resubmit next month."
  },
  "message": "Ticket rejected"
}
```

**Error Response (403 Forbidden):**

```json
{
  "error": "Only approver can approve/reject tickets"
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "This ticket does not require approval"
}
```

---

### 2.6 Add Comment to GA Ticket

**Endpoint:** `POST /api/tickets/:id/comments`

**Description:** Add a comment/update to a ticket

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**

```json
{
  "comment": "Supplies have been delivered to Finance department"
}
```

**Success Response (201 Created):**

```json
{
  "id": 5,
  "ticket_id": 1,
  "user_id": 3,
  "user_name": "Admin GA",
  "user_role": "admin",
  "comment": "Supplies have been delivered to Finance department",
  "created_at": "2025-01-15T15:00:00Z"
}
```

---

### 2.7 Get GA Ticket Comments

**Endpoint:** `GET /api/tickets/:id/comments`

**Description:** Get all comments for a specific ticket

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200 OK):**

```json
{
  "comments": [
    {
      "id": 1,
      "ticket_id": 1,
      "user_id": 3,
      "user_name": "Admin GA",
      "user_role": "admin",
      "comment": "Checking availability in inventory",
      "created_at": "2025-01-15T10:30:00Z"
    },
    {
      "id": 2,
      "ticket_id": 1,
      "user_id": 1,
      "user_name": "John Doe",
      "user_role": "requester",
      "comment": "Thank you! We need it by end of week",
      "created_at": "2025-01-15T11:00:00Z"
    },
    {
      "id": 5,
      "ticket_id": 1,
      "user_id": 3,
      "user_name": "Admin GA",
      "user_role": "admin",
      "comment": "Supplies have been delivered to Finance department",
      "created_at": "2025-01-15T15:00:00Z"
    }
  ],
  "total": 3
}
```

---

## 3. Inventory & Asset APIs (Admin Only)

### 3.1 List Assets

**Endpoint:** `GET /api/assets`

**Description:** Get list of all GA assets and inventory (Admin only)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters (Optional):**

```
?category=Office Furniture    # Filter by category
?condition=good               # Filter by condition
?location=Floor 3             # Filter by location
```

**Success Response (200 OK):**

```json
{
  "assets": [
    {
      "id": 1,
      "asset_code": "GA-FUR-001",
      "name": "Office Chair - Ergonomic",
      "category": "Office Furniture",
      "quantity": 50,
      "unit": "pcs",
      "location": "Floor 3 - GA Storage",
      "condition": "good",
      "last_updated": "2025-01-15T09:00:00Z"
    },
    {
      "id": 2,
      "asset_code": "GA-SUP-001",
      "name": "A4 Paper Box",
      "category": "Office Supplies",
      "quantity": 100,
      "unit": "box",
      "location": "Floor 2 - Storage Room",
      "condition": "good",
      "last_updated": "2025-01-14T14:00:00Z"
    },
    {
      "id": 3,
      "asset_code": "GA-PAN-001",
      "name": "Coffee Arabica 1kg",
      "category": "Pantry Supplies",
      "quantity": 20,
      "unit": "pack",
      "location": "Pantry Storage",
      "condition": "good",
      "last_updated": "2025-01-13T10:00:00Z"
    },
    {
      "id": 4,
      "asset_code": "GA-FAC-001",
      "name": "AC Unit Split 1.5 PK",
      "category": "Facility Equipment",
      "quantity": 1,
      "unit": "unit",
      "location": "Meeting Room A",
      "condition": "needs_maintenance",
      "last_maintenance": "2024-12-01",
      "next_maintenance": "2025-02-01",
      "last_updated": "2025-01-15T08:00:00Z"
    }
  ],
  "total": 4
}
```

**GA Asset Categories:**

- `Office Furniture` - Chairs, desks, cabinets, shelves
- `Office Supplies` - Paper, pens, staplers, files, etc (ATK)
- `Pantry Supplies` - Coffee, tea, sugar, snacks, water
- `Facility Equipment` - AC, lighting, building maintenance equipment
- `Meeting Room Equipment` - Projectors, whiteboards, tables, chairs
- `Cleaning Supplies` - Cleaning agents, mops, brooms

**Error Response (403 Forbidden):**

```json
{
  "error": "Only admin can view assets"
}
```

---

### 3.2 Create/Add Asset

**Endpoint:** `POST /api/assets`

**Description:** Add new asset to GA inventory (Admin only)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Whiteboard Marker",
  "category": "Office Supplies",
  "quantity": 50,
  "unit": "pcs",
  "location": "Floor 2 - Storage Room",
  "condition": "good"
}
```

**Field Validations:**

- `name`: required, max 255 characters
- `category`: required (see categories above)
- `quantity`: required, integer
- `unit`: required (pcs, box, pack, unit, etc)
- `location`: required
- `condition`: required, enum: `good`, `needs_maintenance`, `broken`

**Success Response (201 Created):**

```json
{
  "id": 10,
  "asset_code": "GA-SUP-010",
  "name": "Whiteboard Marker",
  "category": "Office Supplies",
  "quantity": 50,
  "unit": "pcs",
  "location": "Floor 2 - Storage Room",
  "condition": "good",
  "last_updated": "2025-01-15T16:00:00Z",
  "message": "Asset added successfully"
}
```

**Error Response (403 Forbidden):**

```json
{
  "error": "Only admin can add assets"
}
```

---

### 3.3 Update Asset Stock

**Endpoint:** `PUT /api/assets/:id/stock`

**Description:** Update asset quantity (Admin only)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**

```json
{
  "quantity": 45,
  "notes": "Distributed 5 boxes to Finance department (Ticket #GA-2025-0001)"
}
```

**Success Response (200 OK):**

```json
{
  "id": 2,
  "asset_code": "GA-SUP-001",
  "name": "A4 Paper Box",
  "quantity": 45,
  "previous_quantity": 50,
  "unit": "box",
  "notes": "Distributed 5 boxes to Finance department (Ticket #GA-2025-0001)",
  "updated_at": "2025-01-15T16:30:00Z",
  "message": "Stock updated successfully"
}
```

---

---

## GA Service Examples

### Example 1: Office Supplies Request (No Approval)

```json
{
  "title": "Monthly ATK for Marketing team",
  "description": "10 pens, 5 notebooks, 2 staplers, 1 box A4 paper",
  "priority": "medium",
  "category": "office_supplies",
  "estimated_cost": 350000
}
```

**Flow:** pending → in_progress → completed → closed

---

### Example 2: Facility Maintenance (With Approval)

```json
{
  "title": "Replace broken AC in CEO room",
  "description": "AC unit not working, needs replacement",
  "priority": "high",
  "category": "facility_maintenance",
  "estimated_cost": 8000000
}
```

**Flow:** pending → waiting_approval → approved → in_progress → completed → closed

---

### Example 3: Pantry Supplies Restock

```json
{
  "title": "Restock pantry supplies for January",
  "description": "Coffee 5kg, sugar 3kg, tea 100 bags, mineral water 10 gallons",
  "priority": "medium",
  "category": "pantry_supplies",
  "estimated_cost": 450000
}
```

**Flow:** pending → in_progress → completed → closed

---

### Example 4: Meeting Room Setup

```json
{
  "title": "Setup Meeting Room B for client presentation",
  "description": "Need projector, flipchart, 20 chairs arranged theater style",
  "priority": "high",
  "category": "meeting_room",
  "estimated_cost": 0
}
```

**Flow:** pending → in_progress → completed → closed
