# mk-pos-billing

A clean Go boilerplate project based on `mk-go-mre` structure.

## Project Structure

```
mk-pos-billing/
├── cmd/api/              # Application entry point
├── internal/             # Private application code
│   ├── api/              # HTTP layer (handlers, middleware, routes)
│   └── infrastructure/   # Infrastructure concerns (database, cache, logger)
├── pkg/                  # Public shared utilities
│   ├── constants/        # Application constants
│   ├── utils/            # Helper utilities
│   └── validation/       # Custom validators
├── .env.example          # Environment configuration template
├── .gitignore           # Git ignore rules
├── dev.dockerfile       # Development Docker setup
├── prod.dockerfile      # Production Docker setup
└── go.mod               # Go module definition
```

## Getting Started

### Prerequisites
- Go 1.25+
- PostgreSQL
- Redis

### Setup

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Update `.env` with your configuration (database, redis, etc.)

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

5. Test the health endpoint:
   ```bash
   curl http://localhost:3000/health
   ```

## Development

This boilerplate includes:
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Logger**: Zap with file rotation and Sentry integration
- **Middleware**: CORS, Auth (JWT), RBAC
- **Utilities**: Response helpers, pagination, time formatting, JWT helpers

## Add Your Features

To add business logic:
1. Create domain models in `internal/domain/`
2. Create repositories in `internal/domain/repository/`
3. Create services in `internal/service/`
4. Create handlers in `internal/api/handlers/`
5. Register routes in `internal/api/routes/`

## Docker

Development:
```bash
docker build -f dev.dockerfile -t mk-pos-billing:dev .
```

Production:
```bash
docker build -f prod.dockerfile -t mk-pos-billing:prod .
```
