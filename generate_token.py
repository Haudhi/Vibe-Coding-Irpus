#!/usr/bin/env python3
"""
Simple JWT token generator for GA Ticketing System API testing.
Generate tokens for different user roles to test API endpoints.
"""

import jwt
import json
from datetime import datetime, timedelta
import uuid

# JWT Configuration (must match server configuration)
JWT_SECRET = "your-super-secret-jwt-key-change-this-in-production"
JWT_ISSUER = "ga-ticketing-system"
ALGORITHM = "HS256"

def generate_token(user_info):
    """Generate a JWT token for testing"""
    now = datetime.utcnow()
    expires_at = now + timedelta(hours=24)

    payload = {
        "user_id": user_info["id"],
        "employee_id": user_info["employee_id"],
        "name": user_info["name"],
        "email": user_info["email"],
        "role": user_info["role"],
        "department": user_info["department"],
        "jti": str(uuid.uuid4()),  # Unique token ID
        "iss": JWT_ISSUER,
        "sub": user_info["id"],
        "aud": ["ga-ticketing-client"],
        "exp": expires_at,
        "nbf": now,
        "iat": now
    }

    token = jwt.encode(payload, JWT_SECRET, algorithm=ALGORITHM)
    # Ensure token is a string (PyJWT 2.x returns string, older versions return bytes)
    if isinstance(token, bytes):
        token = token.decode('utf-8')
    return token

def main():
    """Generate tokens for different user types"""

    # Sample users for testing
    users = {
        "requester": {
            "id": "req-001",
            "employee_id": "EMP001",
            "name": "John Requester",
            "email": "john.requester@company.com",
            "role": "requester",
            "department": "IT"
        },
        "approver": {
            "id": "app-001",
            "employee_id": "EMP002",
            "name": "Jane Approver",
            "email": "jane.approver@company.com",
            "role": "approver",
            "department": "Finance"
        },
        "admin": {
            "id": "adm-001",
            "employee_id": "EMP003",
            "name": "Admin User",
            "email": "admin@company.com",
            "role": "admin",
            "department": "GA"
        }
    }

    print("GA Ticketing System - JWT Token Generator")
    print("=" * 50)
    print()

    # Generate tokens for each user type
    tokens = {}
    for role, user_info in users.items():
        token = generate_token(user_info)
        tokens[role] = token
        print(f"{role.upper()} USER:")
        print(f"  Name: {user_info['name']}")
        print(f"  Email: {user_info['email']}")
        print(f"  Role: {user_info['role']}")
        print(f"  Token: {token}")
        print()

    # Save tokens to file for easy use
    with open("test_tokens.json", "w") as f:
        json.dump(tokens, f, indent=2)

    print("Tokens saved to: test_tokens.json")
    print()
    print("Usage examples:")
    print('export REQUESTER_TOKEN="' + tokens["requester"] + '"')
    print('export APPROVER_TOKEN="' + tokens["approver"] + '"')
    print('export ADMIN_TOKEN="' + tokens["admin"] + '"')

if __name__ == "__main__":
    # Check if PyJWT is installed
    try:
        import jwt
    except ImportError:
        print("Error: PyJWT is not installed.")
        print("Install it with: pip install PyJWT")
        exit(1)

    main()