# Getting Started

## Overview
Vongga Platform is a modern social network application built with Go, featuring a robust backend architecture and real-time capabilities. This document will guide you through setting up and running the project.

## Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.21 or higher
- Docker and Docker Compose
- Git
- Make (optional, but recommended)

## Installation Steps

1. **Clone the Repository**
   ```bash
   git clone [your-repository-url]
   cd backend-go
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Environment Setup**
   ```bash
   cp .env.example .env
   ```
   Edit the `.env` file with your configuration:
   - Database credentials
   - Redis configuration
   - Firebase credentials
   - API keys and secrets

4. **Start Infrastructure Services**
   ```bash
   docker-compose up -d
   ```
   This will start:
   - PostgreSQL database
   - Redis cache
   - Other required services

5. **Run Database Migrations**
   ```bash
   make migrate-up
   ```

6. **Start the Application**
   ```bash
   go run cmd/api/main.go
   ```
   The server will start on the configured port (default: 8080)

## Development Tools

### Hot Reload (Development Mode)
```bash
air
```

### Running Tests
```bash
make test         # Run all tests
make test-cover   # Run tests with coverage report
```

### API Documentation
- Swagger documentation is available at `/swagger/index.html` when the server is running
- Postman collection is available in `/docs/socialnetwork.postman_collection.json`

## Common Issues and Troubleshooting

1. **Database Connection Issues**
   - Verify PostgreSQL is running: `docker ps`
   - Check database credentials in `.env`
   - Ensure database migrations are up to date

2. **Redis Connection Issues**
   - Verify Redis is running: `docker ps`
   - Check Redis connection string in `.env`

3. **Build Errors**
   - Run `go mod tidy` to resolve dependency issues
   - Ensure Go version 1.21 or higher is installed
