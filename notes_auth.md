# GA Ticketing API Documentation - Competition Version (GA Focus Only)

## Base URL

```
http://localhost:8080/api
```

## Authentication

All endpoints except `/auth/login` require JWT token in header:

```
Authorization: Bearer <jwt_token>
```

---

## User Roles

### 3 Roles in System:

1. **Requester** (Employee) - Can create and view own GA tickets
2. **Approver** - Can approve/reject GA requests that need approval
3. **Admin** - Full access to tickets, assets, inventory, and system management

---

## 1. Authentication APIs

### 1.1 Login

**Endpoint:** `POST /api/auth/login`

**Description:** Authenticate user and get JWT token

**Request Headers:**

```
Content-Type: application/json
```

**Request Body:**

```json
{
  "email": "requester@company.com",
  "password": "password123"
}
```

**Success Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "requester@company.com",
    "name": "John Doe",
    "role": "requester"
  }
}
```

**Test Credentials:**

```
requester@company.com / password123 (role: requester)
approver@company.com / password123 (role: approver)
admin@company.com / password123 (role: admin)
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "Invalid credentials"
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "email is required"
}
```

**Curl Examples:**

```bash
# Login as Requester
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "requester@company.com",
    "password": "password123"
  }'

# Login as Approver
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "approver@company.com",
    "password": "password123"
  }'

# Login as Admin
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@company.com",
    "password": "password123"
  }'
```

**Using the Token:**

```bash
# Export token for easy use (replace with actual token from login response)
export AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Use token in subsequent requests
curl -X GET http://localhost:8080/api/tickets \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json"
```

---

### 1.2 Get Current User

**Endpoint:** `GET /api/auth/me`

**Description:** Get current logged-in user information

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200 OK):**

```json
{
  "id": 1,
  "email": "requester@company.com",
  "name": "John Doe",
  "role": "requester",
  "department": "Finance",
  "created_at": "2025-01-15T10:30:00Z"
}
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "Unauthorized"
}
```

**Curl Example:**

```bash
# Get current user info
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json"
```

---

## 4. User APIs

### 4.1 List Users

**Endpoint:** `GET /api/users`

**Description:** Get list of users

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters (Optional):**

```
?role=approver           # Filter by role
```

**Success Response (200 OK):**

```json
{
  "users": [
    {
      "id": 1,
      "email": "requester@company.com",
      "name": "John Doe",
      "role": "requester",
      "department": "Finance"
    },
    {
      "id": 2,
      "email": "approver@company.com",
      "name": "Jane Approver",
      "role": "approver",
      "department": "Management"
    },
    {
      "id": 3,
      "email": "admin@company.com",
      "name": "Admin GA",
      "role": "admin",
      "department": "General Affairs"
    }
  ],
  "total": 3
}
```

---

## 5. Dashboard Stats API

### 5.1 Get Dashboard Statistics

**Endpoint:** `GET /api/stats`

**Description:** Get GA dashboard statistics (role-based view)

**Request Headers:**

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200 OK) - Admin:**

```json
{
  "tickets": {
    "total": 45,
    "pending": 8,
    "waiting_approval": 5,
    "approved": 3,
    "in_progress": 12,
    "completed": 15,
    "closed": 2,
    "by_priority": {
      "low": 15,
      "medium": 20,
      "high": 10
    },
    "by_category": {
      "office_supplies": 18,
      "facility_maintenance": 10,
      "pantry_supplies": 8,
      "meeting_room": 5,
      "office_furniture": 3,
      "general_service": 1
    },
    "total_cost": 45000000,
    "pending_approval_cost": 15000000
  },
  "assets": {
    "total_items": 250,
    "total_categories": 6,
    "low_stock_items": 8,
    "needs_maintenance": 3,
    "by_category": {
      "Office Furniture": 50,
      "Office Supplies": 120,
      "Pantry Supplies": 40,
      "Facility Equipment": 15,
      "Meeting Room Equipment": 20,
      "Cleaning Supplies": 5
    }
  },
  "recent_activities": [
    {
      "type": "ticket_completed",
      "ticket_number": "GA-2025-0015",
      "title": "Request office stationery supplies",
      "timestamp": "2025-01-15T16:30:00Z"
    },
    {
      "type": "ticket_approved",
      "ticket_number": "GA-2025-0014",
      "title": "AC maintenance in meeting room",
      "timestamp": "2025-01-15T15:00:00Z"
    }
  ]
}
```

**Success Response (200 OK) - Requester:**

```json
{
  "my_tickets": {
    "total": 8,
    "pending": 2,
    "waiting_approval": 1,
    "in_progress": 3,
    "completed": 2,
    "by_category": {
      "office_supplies": 5,
      "facility_maintenance": 2,
      "pantry_supplies": 1
    }
  },
  "recent_tickets": [
    {
      "id": 15,
      "ticket_number": "GA-2025-0015",
      "title": "Request office stationery supplies",
      "status": "completed",
      "created_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

**Success Response (200 OK) - Approver:**

```json
{
  "pending_approvals": {
    "total": 5,
    "total_cost": 15000000,
    "tickets": [
      {
        "id": 20,
        "ticket_number": "GA-2025-0020",
        "title": "Purchase 10 new office chairs",
        "category": "office_furniture",
        "estimated_cost": 15000000,
        "requester_name": "John Doe",
        "created_at": "2025-01-15T14:00:00Z"
      }
    ]
  },
  "approved_this_month": {
    "total": 12,
    "total_cost": 30000000
  },
  "rejected_this_month": {
    "total": 2,
    "total_cost": 5000000
  }
}
```
