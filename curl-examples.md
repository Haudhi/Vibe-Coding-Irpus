# GA Ticketing System - API Curl Examples

This document contains example curl commands for all API endpoints with JWT authentication.

## Configuration

- **Base URL**: `http://localhost:8080`
- **JWT Secret**: `your-super-secret-jwt-key-change-this-in-production` (from .env file)
- **Content-Type**: `application/json`

## Authentication Flow

### 1. Login (Get JWT Token)
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response will contain:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400,
  "user": {
    "id": "user-id",
    "email": "user@example.com",
    "role": "user"
  }
}
```

### 2. Refresh JWT Token
```bash
curl -X POST http://localhost:8080/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

## Ticket Management APIs

### 3. Get All Tickets (with pagination and filtering)
```bash
curl -X GET "http://localhost:8080/v1/tickets?page=1&limit=20&status=open&priority=high" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### 4. Create New Ticket
```bash
curl -X POST http://localhost:8080/v1/tickets \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Network Issue in Office",
    "description": "Internet connection is intermittent in the main office area",
    "priority": "high",
    "category": "network",
    "requested_for": {
      "name": "John Doe",
      "email": "john.doe@company.com",
      "department": "IT"
    },
    "location": {
      "building": "Main Building",
      "floor": "3rd Floor",
      "room": "301"
    }
  }'
```

### 5. Get Ticket by ID
```bash
curl -X GET http://localhost:8080/v1/tickets/ticket-123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### 6. Update Ticket
```bash
curl -X PUT http://localhost:8080/v1/tickets/ticket-123 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Network Issue in Office",
    "description": "Internet connection is intermittent in the main office area - Investigation ongoing",
    "priority": "medium",
    "status": "in_progress"
  }'
```

### 7. Assign Ticket to Admin
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/assign \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "assigned_to": "admin-456",
    "assigned_by": "admin-789",
    "note": "Assigning to network specialist"
  }'
```

### 8. Get Ticket Comments
```bash
curl -X GET http://localhost:8080/v1/tickets/ticket-123/comments \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### 9. Add Comment to Ticket
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/comments \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "content": "I have checked the network switch and found some configuration issues. Working on resolving them.",
    "comment_type": "update",
    "is_internal": false,
    "attachments": [
      {
        "name": "network_config.txt",
        "url": "http://localhost:8080/files/attachments/network_config.txt"
      }
    ]
  }'
```

### 10. Approve Ticket
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/approve \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "approved_by": "manager-456",
    "note": "Approved for immediate action",
    "approved_at": "2024-01-15T10:30:00Z"
  }'
```

### 11. Reject Ticket
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/reject \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "rejected_by": "manager-456",
    "reason": "Not within maintenance scope - please contact facilities",
    "rejected_at": "2024-01-15T10:30:00Z"
  }'
```

## Asset Management APIs (Admin Only)

### 12. Get All Assets
```bash
curl -X GET "http://localhost:8080/v1/assets?page=1&limit=20&category=hardware" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### 13. Create New Asset
```bash
curl -X POST http://localhost:8080/v1/assets \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dell Laptop XPS 15",
    "category": "hardware",
    "type": "laptop",
    "serial_number": "DLXPS15001234",
    "model": "XPS 15 9530",
    "manufacturer": "Dell",
    "purchase_date": "2024-01-10",
    "warranty_expiry": "2027-01-10",
    "status": "available",
    "location": {
      "building": "Main Building",
      "floor": "2nd Floor",
      "room": "IT Storage"
    },
    "specifications": {
      "cpu": "Intel Core i7-13700H",
      "ram": "32GB DDR5",
      "storage": "1TB NVMe SSD",
      "display": "15.6\" FHD+"
    },
    "cost": {
      "purchase_price": 2499.99,
      "currency": "USD"
    }
  }'
```

### 14. Get Asset by ID
```bash
curl -X GET http://localhost:8080/v1/assets/asset-456 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### 15. Update Asset
```bash
curl -X PUT http://localhost:8080/v1/assets/asset-456 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "status": "assigned",
    "assigned_to": "user-789",
    "location": {
      "building": "Main Building",
      "floor": "3rd Floor",
      "room": "301"
    },
    "notes": "Assigned to John Doe for development work"
  }'
```

### 16. Update Asset Inventory
```bash
curl -X POST http://localhost:8080/v1/assets/asset-456/inventory \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 1,
    "action": "check_in",
    "reason": "Asset returned after project completion",
    "performed_by": "admin-123",
    "notes": "Asset in good condition"
  }'
```

## User Management APIs

### 17. Get Current User Profile
```bash
curl -X GET http://localhost:8080/v1/users/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

## Health Check API (No Authentication Required)

### 18. Health Check
```bash
curl -X GET http://localhost:8080/health \
  -H "Content-Type: application/json"
```

## Current Basic APIs (Implemented)

### 19. Basic Tickets Endpoint
```bash
curl -X GET http://localhost:8080/api/tickets \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

## Usage Notes

1. **Replace Tokens**: All `Bearer` tokens need to be replaced with actual JWT tokens obtained from the login endpoint
2. **Admin Required**: Asset management endpoints (`/v1/assets/*`) require admin role
3. **Pagination**: Use `page` and `limit` query parameters for paginated results
4. **Filtering**: Use query parameters like `status`, `priority`, `category` for filtering
5. **Error Handling**: Check HTTP status codes and response messages for errors
6. **CORS**: If testing from different origins, ensure CORS is properly configured

## JWT Token Information

- **Secret Key**: `your-super-secret-jwt-key-change-this-in-production`
- **Algorithm**: HS256
- **Expiry**: 24 hours (configurable via `JWT_EXPIRY`)
- **Refresh Expiry**: 168 hours (7 days, configurable via `JWT_REFRESH_EXPIRY`)

## Common Response Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized (invalid/missing JWT)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error