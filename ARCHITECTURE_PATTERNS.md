# mk-go-mre Architecture & Coding Patterns

This document captures the architecture, layering, and coding conventions used in this repository so the same patterns can be reused in new services.

## Core Architecture

- Transport: Gin HTTP server with middleware for auth, RBAC, and OpenAPI auth.
- DI: Google Wire composes handlers, services, repositories, and infrastructure.
- Data access: GORM repositories with explicit transaction helpers and soft-delete patterns.
- Background jobs: Asynq for queues with a job-execution tracking table and middleware.
- Validation: `pkg/validation` wraps go-playground/validator with typed helpers.
- Responses: standardized success and error shapes in `pkg/utils/response.go`.

## Project Structure

- `cmd/api/main.go`: bootstraps env loading, logger, Gin, routes, and CLI entry.
- `cmd/commands/`: CLI jobs and worker commands.
- `internal/api/`: handlers, request/response DTOs, middleware, and routing.
- `internal/container/`: Wire providers and `App` aggregate.
- `internal/domain/`: domain models and repository contracts/implementations.
- `internal/service/`: business logic and orchestration.
- `internal/infrastructure/`: config loaders, DB/Redis, queues, search, cache, SMS.
- `internal/queueables/`: Asynq task payloads and handlers in a single package.
- `pkg/utils` and `pkg/validation`: reusable helpers.
- `api-flow/`: plain-text request flow docs per endpoint.

## Request Lifecycle (HTTP)

1. `main.go` creates the Gin engine and registers routes.
1. Middleware authenticates and attaches context (`AuthContext`, `RequestContext`).
1. Handlers parse and validate input via `pkg/validation` helpers.
1. Handler calls a service, passing tenant/user context and validated input.
1. Service executes business logic, using repositories and transactions as needed.
1. Handler maps domain errors to HTTP status and uses `RespondJSON`/`RespondError`.

## Handler Pattern

- Handlers are thin. They do not talk to the DB directly.
- Use `requireRequestContext` to enforce middleware context.
- Use `validation.GetValidBodyData` / `GetValidUriData` helpers for typed data.
- Use `utils.RespondJSON` and `utils.RespondError` for uniform API responses.
- Map service errors to HTTP status in a local `mapXError` helper.

## Validation Pattern

- Request structs live in `internal/api/request` with `validate` tags.
- Each request type implements `Messages()` and `Attributes()` for errors.
- Validation is triggered in handlers via `validation.GetValidBodyData` and friends.

## Service Pattern

- Services depend on repositories and infrastructure clients via DI.
- Services own business rules and orchestrate multiple repositories.
- Transactions are used for multi-write operations via `db.Transaction`.
- Services return DTOs (not raw GORM models) for handler responses.

## Repository Pattern

- Each repository encapsulates GORM access for a model or aggregate.
- `Tx(tx *gorm.DB)` pattern selects a transaction or the base DB.
- Repositories provide focused methods (`FindByID`, `ListPaged`, etc.).
- Soft deletes are performed with `gorm.DeletedAt` and audit fields.

## Background Jobs (Asynq)

- Payloads and handlers live in `internal/queueables` to avoid import cycles.
- Jobs are registered via `internal/infrastructure/queue/handlers/default.go`.
- Job tracking uses `job_executions` and `JobTrackingMiddleware`.
- Reliable dispatch uses `JobDispatcher.CreateJobRecord` inside a DB transaction, then `DispatchJob` after commit.

## Dependency Injection (Wire)

- `internal/container/wire.go` defines providers for config, DB, repositories, services, handlers, and queue dependencies.
- `App` struct aggregates all handlers and queue dependencies for bootstrapping.
- Regenerate with `wire` after provider changes.

## Logging & Errors

- Logger is Zap with optional Sentry integration.
- Handler logs include request context and errors for debugging.
- Service errors are translated to HTTP statuses by handlers.

## Adding a New Endpoint

1. Define request/response structs in `internal/api/request` and `internal/api/response`.
1. Add route registration in `internal/api/routes` with proper middleware.
1. Implement handler in `internal/api/handlers`.
1. Implement service in `internal/service`.
1. Add repository methods in `internal/domain/repository`.
1. Add DI providers to `internal/container/wire.go` and regenerate with `wire`.
1. Add `api-flow/` plain-text flow doc (required).
1. Add tests where possible (`*_test.go`).

## Adding a New Background Job

1. Define payload and constructors in `internal/queueables/<domain>_tasks.go`.
1. Implement handler in `internal/queueables/<domain>_handlers.go`.
1. Register handler in `internal/infrastructure/queue/handlers/default.go`.
1. Use `JobDispatcher` for safe dispatch with tracking.

## Operational Notes

- Local dev expects `.env` and Redis via `docker-compose.dev.yml`.
- CLI commands run through `cmd/api/main.go` when arguments are supplied.
- Keep `api-flow/` docs in sync with endpoint behavior.
