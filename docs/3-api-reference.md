# API Reference

## Overview

The Vongga Platform API is a RESTful service that provides social networking functionality. This document outlines the available endpoints, their purposes, and usage examples.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## API Endpoints

### Authentication

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword",
    "name": "User Name"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "securepassword"
}
```

### User Management

#### Get User Profile
```http
GET /users/{userId}
Authorization: Bearer <token>
```

#### Update User Profile
```http
PUT /users/{userId}
Authorization: Bearer <token>
Content-Type: application/json

{
    "name": "Updated Name",
    "bio": "Updated bio"
}
```

### Posts

#### Create Post
```http
POST /posts
Authorization: Bearer <token>
Content-Type: application/json

{
    "content": "Post content",
    "media_urls": ["url1", "url2"]
}
```

#### Get Posts
```http
GET /posts
Authorization: Bearer <token>
Query Parameters:
- page (optional): Page number
- limit (optional): Items per page
```

#### Get Post by ID
```http
GET /posts/{postId}
Authorization: Bearer <token>
```

### Comments

#### Add Comment
```http
POST /posts/{postId}/comments
Authorization: Bearer <token>
Content-Type: application/json

{
    "content": "Comment content"
}
```

#### Get Comments
```http
GET /posts/{postId}/comments
Authorization: Bearer <token>
Query Parameters:
- page (optional): Page number
- limit (optional): Items per page
```

## WebSocket API

### Connection

Connect to WebSocket endpoint:
```
ws://localhost:8080/ws
```

### Events

1. **Message Event**
   ```json
   {
       "type": "message",
       "data": {
           "content": "message content",
           "recipient_id": "user_id"
       }
   }
   ```

2. **Notification Event**
   ```json
   {
       "type": "notification",
       "data": {
           "type": "like",
           "post_id": "post_id"
       }
   }
   ```

## Error Responses

The API uses standard HTTP status codes and returns errors in the following format:

```json
{
    "error": {
        "code": "ERROR_CODE",
        "message": "Error description"
    }
}
```

Common error codes:
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## Rate Limiting

The API implements rate limiting to prevent abuse. Limits are as follows:
- 100 requests per minute for authenticated users
- 20 requests per minute for unauthenticated users

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
```
