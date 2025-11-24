# JWT Implementation Test

## 1. Test User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{
  "nim": "434231027",
  "name": "Aryo Prabowo", 
  "email": "aryo@example.com",
  "password": "password123",
  "role": "mahasiswa"
}'
```

## 2. Test User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{
  "email": "aryo@example.com",
  "password": "password123"
}'
```

## 3. Test Protected Endpoint

```bash
# Get the JWT token from login response and use it
curl -X GET http://localhost:8080/api/v1/users/profile \
-H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## 4. Test Achievement Creation (Mahasiswa only)

```bash
curl -X POST http://localhost:8080/api/v1/achievements \
-H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
-H "Content-Type: application/json" \
-d '{
  "title": "Juara 1 Lomba Programming",
  "description": "Mendapat juara 1 dalam lomba programming tingkat nasional",
  "category": "Programming", 
  "achievement_date": "2025-11-20T00:00:00Z"
}'
```