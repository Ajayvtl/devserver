# DevServer API Endpoints (v1)

This document describes the public RESTful APIs for DevServer configuration and platform management.

## Authentication Flow
DevServer uses standard JWT authentication with Bearer tokens. 

1. **Login**: POST `/api/v1/auth/login` returns an access token and refresh token.
2. **Access**: Include the access token in the `Authorization` header: `Authorization: Bearer <token>`.
3. **Organization Context**: Multi-tenant endpoints require the `X-Org-ID` header to evaluate RBAC permissions.

## Standard Error Response
All endpoints follow this error signature:
```json
{
  "success": false,
  "code": "INVALID_REQUEST",
  "message": "Detailed error message",
  "details": null
}
```

## Standard Success Response
All endpoints follow this success signature:
```json
{
  "success": true,
  "data": { ... }
}
```

## Authorization Matrix
| Endpoint | Method | Context Required | RBAC Required |
|----------|--------|-------------------|---------------|
| `/api/v1/auth/login` | POST | None | None |
| `/api/v1/users/me` | GET | Token | None |
| `/api/v1/organizations` | GET | Token | None |
| `/api/v1/organizations` | POST | Token | None (Creates Org) |
| `/api/v1/settings` | GET/POST | Token, X-Org-ID | `write` |
| `/api/v1/environments`| GET | Token, X-Org-ID | `read` |
| `/api/v1/secrets` | POST | Token, X-Org-ID | `write` |
| `/api/v1/providers` | GET/POST | Token, X-Org-ID | `write` |

## Example Requests

### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "provider": "local",
  "username": "admin",
  "password": "password"
}
```

### Create Organization
```http
POST /api/v1/organizations
Content-Type: application/json
Authorization: Bearer <token>

{
  "name": "Acme Corp"
}
```

### Save Settings
```http
POST /api/v1/settings
Content-Type: application/json
Authorization: Bearer <token>
X-Org-ID: <org-id>

{
  "scope": "organization",
  "ownerId": "<org-id>",
  "key": "theme",
  "value": "dark"
}
```
