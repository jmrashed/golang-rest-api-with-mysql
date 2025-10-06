# API Documentation

## Authentication & Authorization System

This API implements a comprehensive authentication and authorization system with JWT tokens, role-based access control (RBAC), and permission-based access control.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

All protected endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <access_token>
```

## Endpoints

### Public Endpoints (No Authentication Required)

#### POST /register
Register a new user account.

**Request Body:**
```json
{
  "username": "string (required, min: 3, max: 50)",
  "email": "string (required, valid email)",
  "password": "string (required, min: 6)"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "is_active": true,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "roles": [
        {
          "id": 2,
          "name": "user",
          "description": "Regular user with limited access",
          "permissions": [
            {
              "id": 1,
              "name": "read_users",
              "description": "Read user information",
              "resource": "users",
              "action": "read"
            }
          ]
        }
      ]
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

#### POST /login
Authenticate user and receive tokens.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "is_active": true,
      "roles": [...]
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

#### POST /refresh
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "string (required)"
}
```

**Response (200 OK):**
```json
{
  "message": "Token refreshed successfully",
  "data": {
    "user": {...},
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

### Protected Endpoints (Authentication Required)

#### GET /profile
Get current user profile.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "message": "Profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "is_active": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z",
    "roles": [...]
  }
}
```

#### PUT /profile
Update current user profile.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "username": "string (optional)",
  "email": "string (optional, valid email)"
}
```

**Response (200 OK):**
```json
{
  "message": "Profile updated successfully",
  "data": {
    "id": 1,
    "username": "newusername",
    "email": "newemail@example.com",
    "is_active": true,
    "updated_at": "2023-01-01T00:00:00Z",
    "roles": [...]
  }
}
```

#### POST /change-password
Change user password.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "current_password": "string (required)",
  "new_password": "string (required, min: 6)"
}
```

**Response (200 OK):**
```json
{
  "message": "Password changed successfully"
}
```

#### POST /logout
Logout from current session (invalidate refresh token).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request Body:**
```json
{
  "refresh_token": "string (required)"
}
```

**Response (200 OK):**
```json
{
  "message": "Logout successful"
}
```

#### POST /logout-all
Logout from all sessions (invalidate all refresh tokens).

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "message": "Logged out from all devices"
}
```

### Admin Endpoints (Admin Role Required)

#### Base Path: /admin

All admin endpoints require the "admin" role.

### Moderator Endpoints (Moderator or Admin Role Required)

#### Base Path: /moderator

All moderator endpoints require either "moderator" or "admin" role.

### Todo Endpoints (Permission-Based Access)

#### Base Path: /todos

All todo endpoints require the "read_todos" permission.

## Error Responses

All error responses follow this format:

```json
{
  "error": "HTTP Status Text",
  "message": "Detailed error message"
}
```

### Common Error Codes

- **400 Bad Request**: Invalid request data or validation errors
- **401 Unauthorized**: Missing, invalid, or expired authentication token
- **403 Forbidden**: Insufficient permissions or role
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict (e.g., username already exists)
- **500 Internal Server Error**: Server-side error

## Authentication Flow

1. **Registration**: User registers with username, email, and password
2. **Login**: User authenticates and receives access and refresh tokens
3. **API Access**: User includes access token in Authorization header for protected endpoints
4. **Token Refresh**: When access token expires, use refresh token to get new tokens
5. **Logout**: Invalidate refresh tokens when user logs out

## Security Features

- **Password Hashing**: Passwords are hashed using bcrypt
- **JWT Tokens**: Stateless authentication with signed JWT tokens
- **Token Expiration**: Access tokens expire in 15 minutes, refresh tokens in 7 days
- **Role-Based Access Control**: Users have roles that determine access levels
- **Permission-Based Access Control**: Fine-grained permissions for specific actions
- **Secure Token Storage**: Refresh tokens are hashed before database storage
- **CORS Support**: Cross-origin resource sharing enabled

## Roles and Permissions

### Default Roles

1. **admin**: Full system access
   - All permissions

2. **user**: Basic user access
   - read_users
   - read_todos
   - write_todos

3. **moderator**: Intermediate access
   - read_users
   - write_users
   - read_todos
   - write_todos
   - delete_todos

### Available Permissions

- **read_users**: Read user information
- **write_users**: Create and update users
- **delete_users**: Delete users
- **read_todos**: Read todos
- **write_todos**: Create and update todos
- **delete_todos**: Delete todos
- **manage_roles**: Manage user roles

## Rate Limiting

Currently not implemented but recommended for production:
- Login attempts: 5 per minute per IP
- Registration: 3 per hour per IP
- API calls: 1000 per hour per user

## Health Check

#### GET /health
Check API health status.

**Response (200 OK):**
```
OK
```