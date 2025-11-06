# GA Ticketing System - Complete API Curl Examples with JWT Tokens

This document contains curl examples for all API endpoints with working JWT tokens for testing.

## Configuration

- **Base URL**: `http://localhost:8080`
- **JWT Secret**: `your-super-secret-jwt-key-change-this-in-production`
- **Content-Type**: `application/json`

## Authentication Tokens

The following JWT tokens are pre-generated for testing:

### Requester User (John Requester)
```
REQUESTER_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY"
```

### Approver User (Jane Approver)
```
APPROVER_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYXBwLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAyIiwibmFtZSI6IkphbmUgQXBwcm92ZXIiLCJlbWFpbCI6ImphbmUuYXBwcm92ZXJAY29tcGFueS5jb20iLCJyb2xlIjoiYXBwcm92ZXIiLCJkZXBhcnRtZW50IjoiRmluYW5jZSIsImp0aSI6IjljNTNiZGYzLWM5NTUtNDA4Mi1iMTQ1LTNkMjI2ZDM0M2UyMCIsImlzcyI6ImdhLXRpY2tldGluZy1zeXN0ZW0iLCJzdWIiOiJhcHAtMDAxIiwiYXVkIjpbImdhLXRpY2tldGluZy1jbGllbnQiXSwiZXhwIjoxNzYyNTEzMDIyLCJuYmYiOjE3NjI0MjY2MjIsImlhdCI6MTc2MjQyNjYyMn0.9w6vUazPAOk22256qu4yEByuE74ZHnb6WBcGpJAocu4"
```

### Admin User (Admin User)
```
ADMIN_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI"
```

---

## API ENDPOINTS

### 1. Health Check (No Authentication Required)

```bash
curl -X GET http://localhost:8080/health \
  -H "Content-Type: application/json"
```

**Expected Response:**
```json
{
  "status": "healthy",
  "service": "ga-ticketing"
}
```

---

### 2. Simple Mock API Endpoints (No Authentication Required)

#### 2.1 Get All Tickets (Mock)
```bash
curl -X GET http://localhost:8080/api/tickets \
  -H "Content-Type: application/json"
```

#### 2.2 Create Ticket (Mock)
```bash
curl -X POST http://localhost:8080/api/tickets \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Need office supplies",
    "description": "Running out of pens and notebooks",
    "priority": "medium",
    "category": "office_supplies",
    "estimated_cost": 150000,
    "requester_name": "John Doe",
    "department": "IT"
  }'
```

#### 2.3 Get Single Ticket (Mock)
```bash
curl -X GET http://localhost:8080/api/tickets/1 \
  -H "Content-Type: application/json"
```

#### 2.4 Update Ticket Status (Mock)
```bash
curl -X PUT http://localhost:8080/api/tickets/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "assigned_to": "admin",
    "actual_cost": 145000,
    "notes": "Approved and processing"
  }'
```

#### 2.5 Approve/Reject Ticket (Mock)
```bash
curl -X PUT http://localhost:8080/api/tickets/1/approval \
  -H "Content-Type: application/json" \
  -d '{
    "action": "approve",
    "approver_name": "Manager",
    "notes": "Approved for processing"
  }'
```

#### 2.6 Add Comment (Mock)
```bash
curl -X POST http://localhost:8080/api/tickets/1/comments \
  -H "Content-Type: application/json" \
  -d '{
    "comment": "Working on this request now",
    "author_name": "Admin"
  }'
```

#### 2.7 Get Comments (Mock)
```bash
curl -X GET http://localhost:8080/api/tickets/1/comments \
  -H "Content-Type: application/json"
```

#### 2.8 Get All Assets (Mock)
```bash
curl -X GET http://localhost:8080/api/assets \
  -H "Content-Type: application/json"
```

#### 2.9 Create Asset (Mock)
```bash
curl -X POST http://localhost:8080/api/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Office Chair",
    "description": "Ergonomic office chair",
    "category": "office_furniture",
    "quantity": 10,
    "location": "Main Office",
    "condition": "good"
  }'
```

#### 2.10 Update Asset Stock (Mock)
```bash
curl -X PUT http://localhost:8080/api/assets/1/stock \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 15,
    "operation": "add",
    "notes": "New chairs delivered",
    "updated_by": "Admin"
  }'
```

---

### 3. Structured API Endpoints (JWT Authentication Required)

#### 3.1 Get All Tickets (Authenticated)

**Requester Role:**
```bash
curl -X GET "http://localhost:8080/v1/tickets?page=1&limit=20&status=pending" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY" \
  -H "Content-Type: application/json"
```

#### 3.2 Create Ticket (Authenticated)

**Requester Role:**
```bash
curl -X POST http://localhost:8080/v1/tickets \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Request for Office Supplies",
    "description": "Need additional office supplies for the IT department including pens, notebooks, and printer paper",
    "category": "office_supplies",
    "priority": "medium",
    "estimated_cost": 250000,
    "requester_id": "req-001"
  }'
```

#### 3.3 Get Single Ticket (Authenticated)

**Requester Role:**
```bash
curl -X GET http://localhost:8080/v1/tickets/ticket-123 \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY" \
  -H "Content-Type: application/json"
```

#### 3.4 Update Ticket (Authenticated)

**Requester Role:**
```bash
curl -X PUT http://localhost:8080/v1/tickets/ticket-123 \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated: Request for Office Supplies",
    "description": "Need additional office supplies for the IT department including pens, notebooks, and printer paper - Updated quantities",
    "priority": "high",
    "estimated_cost": 300000,
    "reason": "Updated cost due to price increase"
  }'
```

#### 3.5 Assign Ticket (Admin Only)

**Admin Role:**
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/assign \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI" \
  -H "Content-Type: application/json" \
  -d '{
    "admin_id": "adm-001"
  }'
```

#### 3.6 Add Comment to Ticket (Authenticated)

**Requester Role:**
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/comments \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Please prioritize this request as we need the supplies urgently"
  }'
```

#### 3.7 Get Ticket Comments (Authenticated)

**Requester Role:**
```bash
curl -X GET "http://localhost:8080/v1/tickets/ticket-123/comments?page=1&limit=20" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY" \
  -H "Content-Type: application/json"
```

#### 3.8 Approve Ticket (Approver/Admin Only)

**Approver Role:**
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/approval/approve \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYXBwLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAyIiwibmFtZSI6IkphbmUgQXBwcm92ZXIiLCJlbWFpbCI6ImphbmUuYXBwcm92ZXJAY29tcGFueS5jb20iLCJyb2xlIjoiYXBwcm92ZXIiLCJkZXBhcnRtZW50IjoiRmluYW5jZSIsImp0aSI6IjljNTNiZGYzLWM5NTUtNDA4Mi1iMTQ1LTNkMjI2ZDM0M2UyMCIsImlzcyI6ImdhLXRpY2tldGluZy1zeXN0ZW0iLCJzdWIiOiJhcHAtMDAxIiwiYXVkIjpbImdhLXRpY2tldGluZy1jbGllbnQiXSwiZXhwIjoxNzYyNTEzMDIyLCJuYmYiOjE3NjI0MjY2MjIsImlhdCI6MTc2MjQyNjYyMn0.9w6vUazPAOk22256qu4yEByuE74ZHnb6WBcGpJAocu4" \
  -H "Content-Type: application/json" \
  -d '{
    "comments": "Approved. The office supplies are necessary for team operations."
  }'
```

#### 3.9 Reject Ticket (Approver/Admin Only)

**Approver Role:**
```bash
curl -X POST http://localhost:8080/v1/tickets/ticket-123/approval/reject \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYXBwLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAyIiwibmFtZSI6IkphbmUgQXBwcm92ZXIiLCJlbWFpbCI6ImphbmUuYXBwcm92ZXJAY29tcGFueS5jb20iLCJyb2xlIjoiYXBwcm92ZXIiLCJkZXBhcnRtZW50IjoiRmluYW5jZSIsImp0aSI6IjljNTNiZGYzLWM5NTUtNDA4Mi1iMTQ1LTNkMjI2ZDM0M2UyMCIsImlzcyI6ImdhLXRpY2tldGluZy1zeXN0ZW0iLCJzdWIiOiJhcHAtMDAxIiwiYXVkIjpbImdhLXRpY2tldGluZy1jbGllbnQiXSwiZXhwIjoxNzYyNTEzMDIyLCJuYmYiOjE3NjI0MjY2MjIsImlhdCI6MTc2MjQyNjYyMn0.9w6vUazPAOk22256qu4yEByuE74ZHnb6WBcGpJAocu4" \
  -H "Content-Type: application/json" \
  -d '{
    "comments": "Rejected. Budget constraints - please resubmit with reduced quantities."
  }'
```

#### 3.10 Get All Assets (Admin Only)

**Admin Role:**
```bash
curl -X GET "http://localhost:8080/v1/assets?page=1&limit=20&category=office_furniture" \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI" \
  -H "Content-Type: application/json"
```

#### 3.11 Create Asset (Admin Only)

**Admin Role:**
```bash
curl -X POST http://localhost:8080/v1/assets \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ergonomic Office Chair",
    "description": "High-back ergonomic office chair with lumbar support",
    "category": "office_furniture",
    "quantity": 15,
    "location": "Main Office Floor 3",
    "unit_cost": 2500000
  }'
```

#### 3.12 Get Single Asset (Admin Only)

**Admin Role:**
```bash
curl -X GET http://localhost:8080/v1/assets/asset-456 \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI" \
  -H "Content-Type: application/json"
```

#### 3.13 Update Asset (Admin Only)

**Admin Role:**
```bash
curl -X PUT http://localhost:8080/v1/assets/asset-456 \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI" \
  -H "Content-Type: application/json" \
  -d '{
    "location": "Main Office Floor 4",
    "condition": "needs_maintenance",
    "unit_cost": 2750000,
    "description": "High-back ergonomic office chair with lumbar support - Price updated"
  }'
```

#### 3.14 Update Asset Inventory (Admin Only)

**Admin Role:**
```bash
curl -X POST http://localhost:8080/v1/assets/asset-456/inventory \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI" \
  -H "Content-Type: application/json" \
  -d '{
    "change_type": "add",
    "quantity": 5,
    "reason": "New chairs delivered from supplier"
  }'
```

#### 3.15 Get Pending Approvals (Approver/Admin Only)

**Approver Role:**
```bash
curl -X GET http://localhost:8080/v1/approvals/pending \
  -H "Authorization: Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYXBwLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAyIiwibmFtZSI6IkphbmUgQXBwcm92ZXIiLCJlbWFpbCI6ImphbmUuYXBwcm92ZXJAY29tcGFueS5jb20iLCJyb2xlIjoiYXBwcm92ZXIiLCJkZXBhcnRtZW50IjoiRmluYW5jZSIsImp0aSI6IjljNTNiZGYzLWM5NTUtNDA4Mi1iMTQ1LTNkMjI2ZDM0M2UyMCIsImlzcyI6ImdhLXRpY2tldGluZy1zeXN0ZW0iLCJzdWIiOiJhcHAtMDAxIiwiYXVkIjpbImdhLXRpY2tldGluZy1jbGllbnQiXSwiZXhwIjoxNzYyNTEzMDIyLCJuYmYiOjE3NjI0MjY2MjIsImlhdCI6MTc2MjQyNjYyMn0.9w6vUazPAOk22256qu4yEByuE74ZHnb6WBcGpJAocu4" \
  -H "Content-Type: application/json"
```

---

## Quick Setup Commands

```bash
# Export tokens for easy use
export REQUESTER_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoicmVxLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAxIiwibmFtZSI6IkpvaG4gUmVxdWVzdGVyIiwiZW1haWwiOiJqb2huLnJlcXVlc3RlckBjb21wYW55LmNvbSIsInJvbGUiOiJyZXF1ZXN0ZXIiLCJkZXBhcnRtZW50IjoiSVQiLCJqdGkiOiIzMmFlNTM5Mi0xZDFlLTRkZjktOGI1MS05MjJkMjA1YTliMmIiLCJpc3MiOiJnYS10aWNrZXRpbmctc3lzdGVtIiwic3ViIjoicmVxLTAwMSIsImF1ZCI6WyJnYS10aWNrZXRpbmctY2xpZW50Il0sImV4cCI6MTc2MjUxMzAyMiwibmJmIjoxNzYyNDI2NjIyLCJpYXQiOjE3NjI0MjY2MjJ9.qC9nLOOpiNQ6e-LRYCEyZ3OgYuoZXevTnbcCQbB95LY"

export APPROVER_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYXBwLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAyIiwibmFtZSI6IkphbmUgQXBwcm92ZXIiLCJlbWFpbCI6ImphbmUuYXBwcm92ZXJAY29tcGFueS5jb20iLCJyb2xlIjoiYXBwcm92ZXIiLCJkZXBhcnRtZW50IjoiRmluYW5jZSIsImp0aSI6IjljNTNiZGYzLWM5NTUtNDA4Mi1iMTQ1LTNkMjI2ZDM0M2UyMCIsImlzcyI6ImdhLXRpY2tldGluZy1zeXN0ZW0iLCJzdWIiOiJhcHAtMDAxIiwiYXVkIjpbImdhLXRpY2tldGluZy1jbGllbnQiXSwiZXhwIjoxNzYyNTEzMDIyLCJuYmYiOjE3NjI0MjY2MjIsImlhdCI6MTc2MjQyNjYyMn0.9w6vUazPAOk22256qu4yEByuE74ZHnb6WBcGpJAocu4"

export ADMIN_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoiYWRtLTAwMSIsImVtcGxveWVlX2lkIjoiRU1QMDAzIiwibmFtZSI6IkFkbWluIFVzZXIiLCJlbWFpbCI6ImFkbWluQGNvbXBhbnkuY29tIiwicm9sZSI6ImFkbWluIiwiZGVwYXJ0bWVudCI6IkdBIiwianRpIjoiYmMyOThiMjgtMjRhOS00NGY5LWI2MzMtNTE3ZWVjMGJiNjJmIiwiaXNzIjoiZ2EtdGlja2V0aW5nLXN5c3RlbSIsInN1YiI6ImFkbS0wMDEiLCJhdWQiOlsiZ2EtdGlja2V0aW5nLWNsaWVudCJdLCJleHAiOjE3NjI1MTMwMjIsIm5iZiI6MTc2MjQyNjYyMiwiaWF0IjoxNzYyNDI2NjIyfQ.vjQzUXYy-OqF5dtCR67kwV07X_vI7Qvv-bSwVkHlDUI"

# Test health check
curl -X GET http://localhost:8080/health

# Create a ticket with requester token
curl -X POST http://localhost:8080/v1/tickets \
  -H "Authorization: Bearer $REQUESTER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Ticket", "description": "This is a test", "category": "office_supplies", "priority": "medium", "estimated_cost": 100000, "requester_id": "req-001"}'
```

---

## Error Handling

Common HTTP status codes to expect:

- `200` - Success
- `201` - Created successfully
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (missing/invalid JWT)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `500` - Internal Server Error

If you get `401 Unauthorized`, check:
1. Token is correctly set in Authorization header
2. Token hasn't expired
3. JWT_SECRET matches server configuration

If you get `403 Forbidden`, check:
1. User has the required role for the endpoint
2. User has access to the specific resource (e.g., their own tickets)