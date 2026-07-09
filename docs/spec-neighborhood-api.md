# Spec: Neighborhood App REST API

## Objective

Build a complete RESTful CRUD API for a neighborhood community app using Go Gin framework. The API manages users, groups, posts, and notifications based on the ERD defined in `neighborhood_app_erd.mermaid`.

**Users:** Residents who can create groups, join groups, and make posts.
**Groups:** Location-based neighborhood groups with a center point and radius.
**Posts:** Community posts with types (lost_found, give_away, alert, general, service_request).
**Notifications:** Alerts triggered by posts to relevant users.

## Tech Stack

- **Language:** Go 1.25
- **Framework:** Gin (github.com/gin-gonic/gin)
- **Database:** PostgreSQL 17 (via docker-compose)
- **DB Driver:** github.com/lib/pq (or pgx)
- **UUID:** github.com/google/uuid
- **ORM:** No ORM — raw SQL with sqlx or database/sql
- **Auth:** Simple JWT (github.com/golang-jwt/jwt/v5)
- **Validation:** github.com/go-playground/validator/v10
- **Password:** golang.org/x/crypto (bcrypt)

## Commands

```bash
# Run database
docker compose up -d

# Run backend
go run ./cmd/api

# Build
go build ./cmd/api

# Test
go test ./... -v

# Lint
golangci-lint run ./...
```

## Project Structure

```
backend/
├── cmd/api/main.go              → Entry point, server startup
├── config/
│   ├── config.go                → Config struct + loader
│   ├── config.dev.yaml          → Dev config
│   └── config.prod.yaml         → Prod config
├── internal/
│   ├── controller/              → HTTP handlers (Gin controllers)
│   │   ├── user_controller.go
│   │   ├── group_controller.go
│   │   ├── post_controller.go
│   │   └── notification_controller.go
│   ├── dto/                     → Data Transfer Objects (request/response)
│   │   ├── user_dto.go
│   │   ├── group_dto.go
│   │   ├── post_dto.go
│   │   └── notification_dto.go
│   ├── middleware/              → Gin middleware (auth, CORS, logging)
│   │   └── auth_middleware.go
│   ├── model/                   → Database models (structs)
│   │   ├── user.go
│   │   ├── group.go
│   │   ├── post.go
│   │   └── notification.go
│   ├── repository/              → Data access layer (SQL queries)
│   │   ├── user_repository.go
│   │   ├── group_repository.go
│   │   ├── post_repository.go
│   │   └── notification_repository.go
│   ├── routes/                  → Route registration
│   │   └── routes.go
│   └── service/                 → Business logic layer
│       ├── user_service.go
│       ├── group_service.go
│       ├── post_service.go
│       └── notification_service.go
├── migrations/                  → SQL migration files
│   ├── 001_create_users_table.sql
│   ├── 002_create_groups_table.sql
│   ├── 003_create_group_members_table.sql
│   ├── 004_create_posts_table.sql
│   └── 005_create_notifications_table.sql
└── pkg/
    ├── database/
    │   └── postgres.go          → DB connection
    ├── errors/
    │   └── errors.go            → Custom error types
    ├── logger/
    │   └── logger.go            → Logging setup
    ├── response/
    │   └── response.go          → Standard API response helpers
    └── utils/
        └── utils.go             → Shared utilities
```

## API Endpoints

### Users
| Method | Path | Description |
|--------|------|-------------|
| POST   | /api/v1/users | Create user (register) |
| GET    | /api/v1/users | List users (paginated) |
| GET    | /api/v1/users/:id | Get user by ID |
| PUT    | /api/v1/users/:id | Update user |
| DELETE | /api/v1/users/:id | Delete user |
| POST   | /api/v1/auth/login | Login (returns JWT) |

### Groups
| Method | Path | Description |
|--------|------|-------------|
| POST   | /api/v1/groups | Create group |
| GET    | /api/v1/groups | List groups (paginated, filter by location) |
| GET    | /api/v1/groups/:id | Get group by ID |
| PUT    | /api/v1/groups/:id | Update group |
| DELETE | /api/v1/groups/:id | Delete group |
| POST   | /api/v1/groups/:id/join | Join group |
| POST   | /api/v1/groups/:id/leave | Leave group |
| GET    | /api/v1/groups/:id/members | List group members |

### Posts
| Method | Path | Description |
|--------|------|-------------|
| POST   | /api/v1/posts | Create post |
| GET    | /api/v1/posts | List posts (paginated, filter by type/group/location) |
| GET    | /api/v1/posts/:id | Get post by ID |
| PUT    | /api/v1/posts/:id | Update post |
| DELETE | /api/v1/posts/:id | Delete post |
| PATCH  | /api/v1/posts/:id/resolve | Mark post as resolved |

### Notifications
| Method | Path | Description |
|--------|------|-------------|
| GET    | /api/v1/notifications | List notifications for current user |
| PATCH  | /api/v1/notifications/:id/read | Mark notification as read |

## Standard Response Format

```json
// Success
{
  "data": { ... },
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 100,
    "total_pages": 5
  }
}

// Error
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable message",
    "details": {}
  }
}
```

## HTTP Status Codes
- 200: Success (GET, PUT, PATCH)
- 201: Created (POST)
- 204: No Content (DELETE)
- 400: Bad Request (validation error)
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 409: Conflict (duplicate)
- 422: Unprocessable Entity
- 500: Internal Server Error

## Code Style

```go
// Naming: camelCase for vars, PascalCase for exported
// Error handling: always check errors, return early
// SQL: raw queries with parameterized placeholders ($1, $2, ...)

// Controller example
func (uc *UserController) Create(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
        return
    }
    user, err := uc.service.Create(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
        return
    }
    response.Success(c, http.StatusCreated, user)
}
```

## Testing Strategy

- **Unit tests:** Service layer with mocked repositories
- **Integration tests:** Repository layer with test database
- **Test framework:** Go standard `testing` package + `testify/assert`
- **Coverage target:** >70%
- **Test location:** `*_test.go` files alongside source code

## Boundaries

- **Always:** Validate input at controller layer, use parameterized SQL, return consistent error format, run `go test ./...` before commit
- **Ask first:** Database schema changes, adding new dependencies, changing API response format
- **Never:** Commit secrets/hardcoded passwords, use string concatenation in SQL, expose internal errors to client

## Success Criteria

- [ ] All CRUD endpoints for Users, Groups, Posts, Notifications work
- [ ] Pagination works on all list endpoints
- [ ] JWT authentication works (register, login, protected routes)
- [ ] All SQL migrations run successfully
- [ ] `go build ./...` succeeds
- [ ] `go test ./...` passes with >70% coverage

## Open Questions

- None — spec is based on the provided ERD