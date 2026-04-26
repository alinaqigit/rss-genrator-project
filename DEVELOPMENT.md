# Development Guide

Complete guide for developing, testing, and deploying the RSS Generator Project.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Conventions](#project-conventions)
- [Adding Features](#adding-features)
- [Database Development](#database-development)
- [Testing](#testing)
- [Debugging](#debugging)
- [Git Workflow](#git-workflow)
- [Performance Profiling](#performance-profiling)
- [Deployment Checklist](#deployment-checklist)

---

## Development Setup

### Prerequisites

- **Go**: 1.25.6 or later
- **Docker**: Latest version
- **Docker Compose**: Latest version
- **Git**: Latest version
- **PostgreSQL Client** (optional): `psql` for direct database access

### Initial Setup

1. **Clone Repository**

```bash
git clone https://github.com/alinaqigit/rss-generator-project.git
cd rss-generator-project
```

2. **Install Go Dependencies**

```bash
go mod download
go mod verify
```

3. **Create `.env` File**

```bash
cat > .env << EOF
PORT=8080
HOST=127.0.0.1
DATABASE_URL="postgres://postgres:postgres@localhost:5544/rss_generator?sslmode=disable"
EOF
```

4. **Start PostgreSQL**

```bash
docker-compose up -d postgres
```

Wait for database to be healthy:

```bash
docker-compose logs postgres
```

You should see: `database system is ready to accept connections`

5. **Install Tools**

```bash
# goose (migrations)
go install github.com/pressly/goose/v3/cmd/goose@latest

# sqlc (generates db code)
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# golangci-lint (linting)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

6. **Run Migrations**

```bash
goose -dir ./sql/migrations postgres "$DATABASE_URL" up
```

7. **Start Application**

```bash
go run *.go
```

Server should start on `http://localhost:8080`

### VS Code Setup (Recommended)

Install Go extension:

```
ext install golang.go
```

Create `.vscode/settings.json`:

```json
{
  "go.lintOnSave": "package",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "editor.formatOnSave": true,
  "[go]": {
    "editor.defaultFormatter": "golang.go",
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  }
}
```

---

## Project Conventions

### Code Style

1. **File Naming**
   - Handler files: `handler_*.go` (e.g., `handler_user.go`)
   - Middleware files: `middleware*.go` or `middleaare_*.go`
   - Utility files: descriptive names (e.g., `json.go`, `rss.go`)

2. **Function Naming**
   - Exported (public): PascalCase (`CreateUser`)
   - Unexported (private): camelCase (`createUser`)
   - Handler: `handler<Resource>` pattern

3. **Package Organization**
   - `main`: Server setup and routes
   - `internal/auth`: Authentication logic
   - `internal/db`: Database code (auto-generated)
   - `sql/`: Migrations and queries

4. **Comments**
   - Comment all exported functions
   - Explain the "why", not the "what"
   - Use proper English

### Example

```go
// CreateUser registers a new user account and returns it with an API key.
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Formatting

Run Go formatter automatically:

```bash
gofmt -s -w .
```

Or go imports:

```bash
goimports -w .
```

---

## Adding Features

### 1. Adding a New Endpoint

#### Step 1: Create SQL Queries

File: `sql/queries/xxx.sql`

```sql
-- name: GetXXXByID :one
SELECT * FROM xxx WHERE id = $1;

-- name: CreateXXX :one
INSERT INTO xxx (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;
```

#### Step 2: Generate Database Code

```bash
sqlc generate
```

This creates functions in `internal/db/xxx.sql.go`

#### Step 3: Create Handler

File: `handler_xxx.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/alinaqigit/rss-generator-project/internal/db"
)

func (apiCfg *apiConfig) handlerGetXXX(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Call database
    // Return response
    responseWithJson(w, 200, result)
}
```

#### Step 4: Register Route

File: `main.go`

```go
v1router.Get("/xxx", apiCfg.handlerGetXXX)
v1router.Post("/xxx", apiCfg.middlewareAuth(apiCfg.handlerCreateXXX))
```

#### Step 5: Test Endpoint

```bash
curl http://localhost:8080/v1/xxx
```

---

## Database Development

### Working with Migrations

#### Create Migration

```bash
goose -dir ./sql/migrations postgres "$DATABASE_URL" create migration_name sql
```

Creates file: `sql/migrations/NNN_migration_name.sql`

#### Structure

```sql
-- +goose up
-- CREATE TABLE / ALTER TABLE / CREATE INDEX

-- +goose down
-- DROP TABLE / ALTER TABLE / DROP INDEX
```

#### Apply Migrations

```bash
# Apply all pending migrations
goose -dir ./sql/migrations postgres "$DATABASE_URL" up

# Rollback last migration
goose -dir ./sql/migrations postgres "$DATABASE_URL" down

# Check status
goose -dir ./sql/migrations postgres "$DATABASE_URL" status
```

### Writing SQL Queries

File: `sql/queries/xxx.sql`

```sql
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: CreateUser :one
INSERT INTO users (id, username, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $1, updated_at = $2
WHERE id = $3
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

**Query Types:**

| Type    | Returns       | Usage                        |
| ------- | ------------- | ---------------------------- |
| `:one`  | Single row    | Get one record               |
| `:many` | Multiple rows | Get multiple records         |
| `:exec` | No rows       | Modifications without return |

### Direct Database Access

```bash
# Connect to database
psql $DATABASE_URL

# Show tables
\dt

# Show table structure
\d users

# Run query
SELECT * FROM users;
```

---

## Testing

### Unit Tests

Create test file: `handler_user_test.go`

```go
package main

import (
    "testing"
)

func TestHandlerCreateUser(t *testing.T) {
    // Test implementation
}
```

Run tests:

```bash
go test ./... -v
go test -run TestHandlerCreateUser -v
```

### Integration Tests

```bash
# Start clean database
docker-compose down
docker-compose up -d postgres
goose -dir ./sql/migrations postgres "$DATABASE_URL" up

# Run integration tests
go test ./... -tags integration -v
```

### Manual Testing with curl

```bash
# 1. Create user
API_KEY=$(curl -X POST http://localhost:8080/v1/user \
  -H "Content-Type: application/json" \
  -d '{"name": "testuser"}' | jq -r '.api_key')

echo "API Key: $API_KEY"

# 2. Get user
curl -H "Authorization: ApiKey $API_KEY" \
  http://localhost:8080/v1/user

# 3. Get all feeds
curl http://localhost:8080/v1/feed

# 4. Create feed
FEED_ID=$(curl -X POST http://localhost:8080/v1/feed \
  -H "Authorization: ApiKey $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Feed", "url": "https://example.com/feed"}' \
  | jq -r '.id')

# 5. Follow feed
curl -X POST http://localhost:8080/v1/feed-follow \
  -H "Authorization: ApiKey $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"feed_id\": \"$FEED_ID\"}"

# 6. Get posts
curl -H "Authorization: ApiKey $API_KEY" \
  http://localhost:8080/v1/posts
```

### Load Testing

Using `ab` (Apache Bench):

```bash
# 1000 requests, 10 concurrent
ab -n 1000 -c 10 http://localhost:8080/v1/healthz
```

Using `hey`:

```bash
go install github.com/rakyll/hey@latest

# 1000 requests, 10 concurrent
hey -n 1000 -c 10 http://localhost:8080/v1/healthz
```

---

## Debugging

### Enabling Debug Logging

Add to code:

```go
log.Println("Debug:", variable)
```

Run with log output:

```bash
go run *.go 2>&1 | tee debug.log
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debugging
dlv debug

# In debugger:
(dlv) break main.main
(dlv) continue
(dlv) next
(dlv) print variable
(dlv) quit
```

### Database Debugging

```bash
# Check database state
psql $DATABASE_URL

# View all users
SELECT id, username, api_key FROM users;

# Check feeds for user
SELECT f.* FROM feeds f
  JOIN users u ON f.user_id = u.id
  WHERE u.username = 'alice';

# Check recent posts
SELECT id, title, published_at, feed_id FROM posts
  ORDER BY published_at DESC LIMIT 10;

# Check feed follows
SELECT uf.*, f.name FROM feed_follows uf
  JOIN feeds f ON uf.feed_id = f.id
  JOIN users u ON uf.user_id = u.id
  WHERE u.username = 'alice';
```

### Monitoring Background Scraper

Check logs:

```bash
go run *.go 2>&1 | grep -i scrap
```

Output shows:

- Scraping goroutines starting
- Feed fetch successes/failures
- Post creation status
- Timestamps of scrapes

---

## Git Workflow

### Creating Feature Branch

```bash
git checkout -b feature/user-preferences
```

### Before Committing

```bash
# Format code
gofmt -s -w .

# Lint code
golangci-lint run

# Run tests
go test ./... -v

# Run application
go run *.go &
# Test endpoints
curl http://localhost:8080/v1/healthz
```

### Committing Changes

```bash
git add .
git commit -m "feat: add user preference endpoints"
```

**Commit message format:**

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `refactor:` Code refactoring
- `test:` Test additions
- `perf:` Performance improvements

### Pushing Changes

```bash
git push origin feature/user-preferences
```

### Creating Pull Request

On GitHub - describe changes, reference issues, add test results

---

## Performance Profiling

### CPU Profiling

```bash
# Build with pprof
go build -o bin/rss-generator

# Start application with profiling
./bin/rss-generator &

# Collect CPU profile (30 seconds)
curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof

# Analyze
go tool pprof cpu.prof
```

In pprof interactive mode:

```
(pprof) top    # Show top functions by CPU usage
(pprof) list functionName  # Show function source
(pprof) web    # Generate graph (needs graphviz)
```

### Memory Profiling

```bash
# Collect memory profile
curl http://localhost:6060/debug/pprof/heap > mem.prof

# Analyze
go tool pprof mem.prof
```

### Database Query Performance

Add logging wrapper (in production use APM):

```go
start := time.Now()
result, err := apiCfg.DB.GetPosts(ctx, params)
duration := time.Since(start)
if duration > 100*time.Millisecond {
    log.Printf("SLOW QUERY: GetPosts took %v", duration)
}
```

---

## Deployment Checklist

### Pre-Deployment

- [ ] All tests passing: `go test ./... -v`
- [ ] No linting errors: `golangci-lint run`
- [ ] Code formatted: `gofmt -s -w .`
- [ ] Migrations tested on clean database
- [ ] Sensitive values in environment variables (not hardcoded)
- [ ] CORS settings appropriate for environment
- [ ] Database backups configured
- [ ] Monitoring/logging configured

### Environment Configuration

**Development:**

```env
PORT=8080
HOST=127.0.0.1
DATABASE_URL=postgres://postgres:postgres@localhost:5544/rss_generator?sslmode=disable
```

**Staging/Production:**

```env
PORT=8080
HOST=0.0.0.0
DATABASE_URL=postgres://user:pass@prod-db.cloud.com:5432/rss_generator?sslmode=require
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=https://app.example.com,https://www.app.example.com
```

### Docker Deployment

1. **Build Image**

```bash
docker build -t rss-generator:v1.0.0 .
```

2. **Tag for Registry**

```bash
docker tag rss-generator:v1.0.0 myregistry/rss-generator:v1.0.0
docker tag rss-generator:v1.0.0 myregistry/rss-generator:latest
```

3. **Push to Registry**

```bash
docker push myregistry/rss-generator:v1.0.0
docker push myregistry/rss-generator:latest
```

4. **Deploy**

```bash
docker run -d \
  -p 8080:8080 \
  -e PORT=8080 \
  -e HOST=0.0.0.0 \
  -e DATABASE_URL="..." \
  --name rss-generator \
  myregistry/rss-generator:latest
```

### Kubernetes Deployment

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl rollout status deployment/rss-generator
```

### Health Check After Deployment

```bash
# Check health endpoint
curl http://your-domain:8080/v1/healthz

# Check logs
docker logs rss-generator  # or kubectl logs deployment/rss-generator

# Test API
curl http://your-domain:8080/v1/feed
```

### Rollback Plan

If deployment fails:

```bash
# Docker
docker run -d ... myregistry/rss-generator:previous-version

# Kubernetes
kubectl rollout undo deployment/rss-generator
```

---

## Common Development Tasks

### Checking Code Quality

```bash
# Run all checks
./scripts/check-quality.sh  # Create this script

# Manual checks:
golangci-lint run
go vet ./...
go test ./... -v -cover
```

### Updating Dependencies

```bash
# Check for updates
go list -u -m all

# Get specific version
go get github.com/go-chi/chi@v1.5.5

# Update all
go get -u ./...

# Verify
go mod verify
```

### Adding New Dependencies

```bash
go get github.com/new/package@v1.2.3
go mod verify
go mod tidy
git add go.mod go.sum
```

---

**Last Updated:** April 26, 2025
