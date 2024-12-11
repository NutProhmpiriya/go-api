# Social Network Application

A modern social network application built with Go, PostgreSQL, Redis, and AlpineJS.

## Tech Stack

### Backend
- Go (Golang)
- PostgreSQL (Primary Database)
- Redis (Caching)
- Firebase Authentication
- RESTful API (with GraphQL consideration)
- WebSocket support

### Frontend
- AlpineJS
- TailwindCSS

## Project Structure
```
.
├── cmd/                    # Application entry points
├── internal/              # Private application and library code
│   ├── domain/           # Enterprise business rules
│   ├── usecase/          # Application business rules
│   ├── repository/       # Data layer implementations
│   └── delivery/         # Delivery mechanisms (HTTP, WebSocket)
├── pkg/                  # Public library code
├── config/              # Configuration files
├── scripts/             # Build and deployment scripts
└── docker/              # Docker configurations
```

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL
- Redis

### Setup
1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Copy `.env.example` to `.env` and configure your environment variables
4. Run `docker-compose up -d` to start PostgreSQL and Redis
5. Run `go run cmd/api/main.go` to start the server

## Development

### Database Migrations
```bash
make migrate-up    # Apply migrations
make migrate-down  # Rollback migrations
```

### Running Tests
```bash
make test         # Run all tests
make test-cover   # Run tests with coverage
```
