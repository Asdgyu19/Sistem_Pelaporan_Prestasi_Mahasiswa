# JWT Refresh Token Implementation

## Overview
Successfully implemented a comprehensive JWT refresh token system for enhanced security and better user experience in the Student Achievement Reporting System.

## Implementation Details

### 1. Enhanced JWT Utilities (`utils/jwt.go`)
- **Dual Token System**: Generates both access tokens (15 minutes) and refresh tokens (7 days)
- **Token Types**: Access tokens for API requests, refresh tokens for obtaining new access tokens
- **Security Features**: 
  - SHA256 token hashing for database storage
  - Token type validation (access vs refresh)
  - Secure random token generation
  - Proper token header extraction

### 2. Refresh Token Service (`service/refresh_token_service.go`)
- **Database Management**: Store, validate, and revoke refresh tokens
- **Security Tracking**: IP address and User-Agent logging for security monitoring  
- **Token Rotation**: Automatic old token revocation when generating new tokens
- **Cleanup Utilities**: Functions to remove expired tokens and revoke all user tokens

### 3. Enhanced Authentication (`helper/auth_helper.go`)
- **Updated Login**: Now generates token pairs instead of single tokens
- **New Endpoints**:
  - `POST /api/v1/auth/refresh` - Refresh access token
  - `POST /api/v1/auth/revoke` - Revoke specific refresh token
  - `POST /api/v1/auth/logout-all` - Logout from all devices
  - `GET /api/v1/auth/tokens` - List active refresh tokens

### 4. Database Schema (`database/migrations/003_add_refresh_tokens.sql`)
```sql
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP,
    ip_address VARCHAR(45),
    user_agent TEXT,
    CONSTRAINT fk_refresh_tokens_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## API Endpoints

### Authentication Endpoints

#### 1. Login (Enhanced)
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "message": "Login successful",
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "mahasiswa"
  }
}
```

#### 2. Refresh Token
```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

#### 3. Revoke Refresh Token
```
POST /api/v1/auth/revoke
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response:
{
  "message": "Refresh token revoked successfully",
  "status": "success"
}
```

#### 4. Logout All Devices
```
POST /api/v1/auth/logout-all
Authorization: Bearer <access_token>

Response:
{
  "message": "Logged out from all devices successfully",
  "status": "success"
}
```

#### 5. List Active Tokens
```
GET /api/v1/auth/tokens
Authorization: Bearer <access_token>

Response:
{
  "active_tokens": [
    {
      "id": 1,
      "created_at": "2025-11-28T08:00:00Z",
      "last_used_at": "2025-11-28T08:30:00Z",
      "expires_at": "2025-12-05T08:00:00Z",
      "ip_address": "192.168.1.100",
      "user_agent": "Mozilla/5.0..."
    }
  ],
  "total": 1
}
```

## Security Features

### 1. Token Security
- **Short-lived Access Tokens**: 15 minutes to limit exposure if compromised
- **Long-lived Refresh Tokens**: 7 days for convenience with secure storage
- **Token Hashing**: Only SHA256 hashes stored in database, not actual tokens
- **Token Rotation**: Old refresh tokens revoked when new ones are issued

### 2. Tracking and Monitoring  
- **IP Address Logging**: Track where tokens are being used
- **User Agent Tracking**: Identify different devices/browsers
- **Last Used Timestamps**: Monitor token activity
- **Revocation Logging**: Track when and why tokens were revoked

### 3. Access Control
- **Token Type Validation**: Ensure correct token types for different operations
- **User-specific Revocation**: Can revoke all tokens for specific users
- **Device-specific Logout**: Can logout individual devices
- **Automatic Cleanup**: Remove expired tokens to maintain database performance

## Client Implementation Guide

### 1. Token Storage
```javascript
// Store tokens securely
localStorage.setItem('access_token', response.access_token);
localStorage.setItem('refresh_token', response.refresh_token);
```

### 2. API Request Interceptor
```javascript
// Add access token to requests
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### 3. Token Refresh Logic
```javascript
// Handle token refresh on 401 responses
axios.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = localStorage.getItem('refresh_token');
      if (refreshToken) {
        try {
          const response = await axios.post('/api/v1/auth/refresh', {
            refresh_token: refreshToken
          });
          
          localStorage.setItem('access_token', response.data.access_token);
          localStorage.setItem('refresh_token', response.data.refresh_token);
          
          // Retry original request
          error.config.headers.Authorization = `Bearer ${response.data.access_token}`;
          return axios.request(error.config);
        } catch (refreshError) {
          // Refresh failed, redirect to login
          localStorage.clear();
          window.location.href = '/login';
        }
      }
    }
    return Promise.reject(error);
  }
);
```

## Database Setup

### Manual Database Setup (Due to Permission Issues)
Since automatic migrations are blocked by database permissions, execute the following SQL manually:

1. **Connect to your PostgreSQL database**
2. **Execute the SQL from**: `database/migrations/MANUAL_003_add_refresh_tokens.sql`
3. **Verify table creation**: `SELECT * FROM refresh_tokens LIMIT 0;`

### Cleanup Function
The migration includes a cleanup function that can be run periodically:
```sql
SELECT cleanup_expired_refresh_tokens();
```

## Testing

### 1. Test Login with New Response Format
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### 2. Test Token Refresh
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### 3. Test Token Management
```bash
# List active tokens
curl -X GET http://localhost:8080/api/v1/auth/tokens \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Logout all devices
curl -X POST http://localhost:8080/api/v1/auth/logout-all \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Benefits of This Implementation

1. **Enhanced Security**: Short-lived access tokens reduce attack surface
2. **Better UX**: Users stay logged in longer without frequent re-authentication
3. **Multi-device Support**: Users can manage sessions across devices
4. **Security Monitoring**: Track and audit token usage patterns
5. **Graceful Token Handling**: Automatic token refresh in client applications
6. **Administrative Control**: Admins can revoke access for specific users/devices

## Next Steps

1. **Execute Manual Database Migration**: Run the SQL from `MANUAL_003_add_refresh_tokens.sql`
2. **Test All Endpoints**: Verify refresh token functionality works correctly  
3. **Update Client Applications**: Implement token refresh logic in frontend
4. **Monitor Usage**: Set up logging and monitoring for token operations
5. **Security Audit**: Review token lifetimes and security policies

The refresh token system is now fully implemented and ready for use once the database table is created manually!