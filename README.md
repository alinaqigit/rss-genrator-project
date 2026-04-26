# RSS Generator Project

A RESTful API service that aggregates RSS feeds, manages users, tracks feed
subscriptions, and automatically scrapes and stores feed content. Built with Go,
PostgreSQL, and deployed with Docker.

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Architecture](#-architecture)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Database Schema](#-database-schema)
- [API Documentation](#-api-documentation)
- [Installation & Setup](#-installation--setup)
- [Environment Variables](#-environment-variables)
- [Running the Project](#-running-the-project)
- [Development](#-development)
- [Deployment](#-deployment)

---

## 📌 Overview

The RSS Generator Project is a backend service designed to:

- **Manage Users**: Create and authenticate users with API keys
- **Track RSS Feeds**: Users can add RSS feed sources they want to follow
- **Subscribe to Feeds**: Users can follow specific feeds they're interested in
- **Aggregate Content**: Automatically scrapes RSS feeds and stores posts
- **Serve Content**: Provides endpoints for users to retrieve their posts from
  followed feeds

This is a perfect solution for building RSS readers, content aggregation
platforms, or notification systems.

---

## 🎯 Features

✅ **User Management**

- User registration with unique username
- Automatic API key generation for authentication
- User deactivation/deletion

✅ **Feed Management**

- Create and manage RSS feed sources
- View all available feeds in the system
- Filter feeds by user

✅ **Feed Subscriptions**

- Users can follow/subscribe to specific feeds
- Retrieve list of feeds a user is following
- Unsubscribe from feeds

✅ **Content Aggregation**

- Automatic background RSS feed scraping
- Concurrent scraping using goroutines (configurable)
- Parse and store RSS feed items as posts
- Track feed fetch timestamps

✅ **Post Retrieval**

- Fetch posts from user's subscribed feeds
- Pagination support (default limit: 10 posts)

✅ **Security**

- API Key-based authentication
- Authorization middleware for protected endpoints
- CORS support for cross-origin requests

---

## 🏗️ Architecture

### System Overview

```
┌─────────────────────────────────────────────────┐
│              Client / Application               │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │  HTTP API Server     │
        │  (Chi Router)        │
        └──────────────────────┘
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
  ┌─────────────┐    ┌──────────────────┐
  │  API Routes │    │  Auth Middleware │
  │  (v1)       │    │  (API Key based) │
  └─────────────┘    └──────────────────┘
        │
        ▼
  ┌────────────────────────────┐
  │   Database Layer (sqlc)    │
  │   - Users                  │
  │   - Feeds                  │
  │   - Feed Follows           │
  │   - Posts                  │
  └────────────────────────────┘
        │
        ▼
  ┌────────────────────────────┐
  │   PostgreSQL Database      │
  │   (Docker Container)       │
  └────────────────────────────┘

┌─────────────────────────────────────────────────┐
│        Background RSS Scraper (Goroutines)      │
│  - Concurrently fetches feeds                   │
│  - Parses RSS XML                               │
│  - Stores posts in database                     │
└─────────────────────────────────────────────────┘
```

### Request Flow

1. Client sends HTTP request with API key in `Authorization` header
2. Request hits Chi router
3. CORS middleware processes the request
4. Route-specific handler is invoked
5. Auth middleware validates API key (if protected endpoint)
6. Handler executes database queries via sqlc
7. Response is marshaled to JSON and sent back

### Background Scraping Flow

1. Ticker triggers at configured intervals (default: 1 minute)
2. Get next batch of feeds to be fetched (based on concurrency limit)
3. Launch goroutines for each feed
4. For each feed:
   - Fetch RSS feed from URL
   - Parse RSS XML structure
   - For each item in RSS:
     - Create Post record in database
   - Update feed's `last_fetched_at` timestamp

---

## 🛠️ Tech Stack

| Component                 | Technology             | Version   |
| ------------------------- | ---------------------- | --------- |
| **Language**              | Go                     | 1.25.6    |
| **Web Framework**         | Chi Router             | v1.5.5    |
| **Database**              | PostgreSQL             | 16-Alpine |
| **ORM/Query Builder**     | sqlc                   | -         |
| **Authentication**        | API Key (Bearer Token) | -         |
| **CORS**                  | chi/cors               | v1.2.2    |
| **Database Driver**       | lib/pq                 | v1.12.1   |
| **UUID Generation**       | google/uuid            | v1.6.0    |
| **Environment Variables** | godotenv               | v1.5.1    |
| **Containerization**      | Docker                 | -         |
| **Orchestration**         | Docker Compose         | 3.9       |

---

## 📁 Project Structure

```
rss-generator-project/
├── main.go                           # Entry point, server setup
├── models.go                         # Data models and type conversions
├── handler_user.go                   # User endpoint handlers
├── handler_feed.go                   # Feed endpoint handlers
├── handler_feedFollow.go             # Feed follow endpoint handlers
├── handler-readiness.go              # Health check endpoint
├── handler_err.go                    # Error handler demo
├── middleaare_auth.go                # API key authentication middleware
├── json.go                           # JSON response utilities
├── rss.go                            # RSS XML parsing logic
├── scraper.go                        # Background RSS feed scraper
├── sqlc.yaml                         # sqlc configuration
├── docker-compose.yml                # Docker Compose setup
├── go.mod                            # Go module dependencies
├── go.sum                            # Go module checksums
│
├── internal/
│   ├── auth/
│   │   └── auth.go                   # API key extraction from headers
│   └── db/
│       ├── db.go                     # Database connection and initialization
│       ├── models.go                 # Database models (auto-generated by sqlc)
│       ├── users.sql.go              # User queries (auto-generated)
│       ├── feeds.sql.go              # Feed queries (auto-generated)
│       ├── feed_follows.sql.go       # Feed follow queries (auto-generated)
│       └── posts.sql.go              # Post queries (auto-generated)
│
├── sql/
│   ├── migrations/
│   │   ├── 001_users.sql             # Create users table
│   │   ├── 002_users_apikey.sql      # Add API key column
│   │   ├── 003_feeds.sql             # Create feeds table
│   │   ├── 004_feed_follows.sql      # Create feed_follows table
│   │   ├── 005_lastfetchedat.sql     # Add last_fetched_at to feeds
│   │   └── 006_posts.sql             # Create posts table
│   └── queries/
│       ├── users.sql                 # User SQL queries
│       ├── feeds.sql                 # Feed SQL queries
│       ├── feed_follows.sql          # Feed follow SQL queries
│       └── posts.sql                 # Post SQL queries
│
└── tmp/
    └── main                          # Compiled binary (git-ignored)
```

---

## 🗄️ Database Schema

### Users Table

Stores user account information with unique API keys.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL UNIQUE,
    api_key VARCHAR(255) UNIQUE NOT NULL DEFAULT (encode(sha256(random()::text::bytea), 'hex'))
);
```

**Fields:**

- `id`: Unique identifier (UUID)
- `created_at`: Account creation timestamp
- `updated_at`: Last update timestamp
- `username`: Unique username
- `api_key`: Auto-generated SHA256 hex-encoded API key for authentication

---

### Feeds Table

Contains RSS feed sources that users can subscribe to.

```sql
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_fetched_at TIMESTAMPTZ
);
```

**Fields:**

- `id`: Unique feed identifier (UUID)
- `created_at`: Feed creation timestamp
- `updated_at`: Last modification timestamp
- `name`: Human-readable feed name
- `url`: RSS feed URL
- `user_id`: Creator of this feed (foreign key)
- `last_fetched_at`: Timestamp of last RSS scrape

---

### Feed Follows Table

Tracks which users follow which feeds (many-to-many relationship).

```sql
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE (user_id, feed_id)
);
```

**Fields:**

- `id`: Unique follow relationship ID
- `created_at`: When user started following
- `updated_at`: Last update
- `user_id`: User ID (foreign key)
- `feed_id`: Feed being followed (foreign key)
- Constraint: Each user can follow each feed only once

---

### Posts Table

Stores individual articles/items from RSS feeds.

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    url TEXT NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);
```

**Fields:**

- `id`: Unique post identifier
- `created_at`: When post was added to database
- `updated_at`: Last modification timestamp
- `title`: Post/article title
- `description`: Post content/summary (nullable)
- `published_at`: Publication date from RSS feed
- `url`: Unique URL to the post
- `feed_id`: Which feed this post came from

---

## 🔌 API Documentation

### Base URL

```
http://localhost:{PORT}/v1
```

### Authentication

All authenticated endpoints require an `Authorization` header:

```
Authorization: ApiKey {api_key}
```

The API key is automatically generated when a user is created.

---

### Endpoints

#### 1. **Health Check**

Check if the server is running.

```
GET /healthz
```

**Response:**

```json
{}
```

**Status Code:** `200 OK`

---

#### 2. **Create User**

Register a new user and receive an API key.

```
POST /user
Content-Type: application/json

{
  "name": "john_doe"
}
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "john_doe",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

**Status Code:** `201 Created`

**Error Response:**

```json
{
  "error": "Couldn't create user: username already exists"
}
```

**Status Code:** `400 Bad Request`

---

#### 3. **Get Current User**

Retrieve authenticated user's information.

```
GET /user
Authorization: ApiKey {api_key}
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:30:00Z",
  "updated_at": "2025-04-26T10:30:00Z",
  "name": "john_doe",
  "api_key": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"
}
```

**Status Code:** `200 OK`

**Error Response:**

```json
{
  "error": "Auth error: No Authentication info found"
}
```

**Status Code:** `403 Forbidden`

---

#### 4. **Delete User**

Deactivate and delete a user account.

```
DELETE /user
Authorization: ApiKey {api_key}
```

**Response:**

```json
{}
```

**Status Code:** `204 No Content`

---

#### 5. **Create Feed**

Add a new RSS feed source.

```
POST /feed
Authorization: ApiKey {api_key}
Content-Type: application/json

{
  "name": "TechNews Daily",
  "url": "https://example.com/feed/rss"
}
```

**Response:**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:35:00Z",
  "updated_at": "2025-04-26T10:35:00Z",
  "name": "TechNews Daily",
  "url": "https://example.com/feed/rss",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Status Code:** `201 Created`

---

#### 6. **Get All Feeds**

Retrieve all available feeds in the system.

```
GET /feed
```

**Response:**

```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T10:35:00Z",
    "updated_at": "2025-04-26T10:35:00Z",
    "name": "TechNews Daily",
    "url": "https://example.com/feed/rss",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T10:40:00Z",
    "updated_at": "2025-04-26T10:40:00Z",
    "name": "Dev Blog",
    "url": "https://devblog.example.com/feed",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
]
```

**Status Code:** `200 OK`

---

#### 7. **Create Feed Follow**

Subscribe to a feed.

```
POST /feed-follow
Authorization: ApiKey {api_key}
Content-Type: application/json

{
  "feed_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-04-26T10:45:00Z",
  "updated_at": "2025-04-26T10:45:00Z",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "feed_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Status Code:** `201 Created`

---

#### 8. **Get Feed Follows**

List all feeds a user is following.

```
GET /feed-follow
Authorization: ApiKey {api_key}
```

**Response:**

```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T10:45:00Z",
    "updated_at": "2025-04-26T10:45:00Z",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "feed_id": "660e8400-e29b-41d4-a716-446655440000"
  }
]
```

**Status Code:** `200 OK`

---

#### 9. **Delete Feed Follow**

Unsubscribe from a feed.

```
DELETE /feed-follow/{feed-follow-id}
Authorization: ApiKey {api_key}
```

**Path Parameters:**

- `feed-follow-id`: UUID of the feed follow relationship

**Response:**

```json
{}
```

**Status Code:** `204 No Content`

---

#### 10. **Get User Posts**

Retrieve posts from all feeds the user is following (paginated, limit 10).

```
GET /posts
Authorization: ApiKey {api_key}
```

**Response:**

```json
[
  {
    "id": "990e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-04-26T11:00:00Z",
    "updated_at": "2025-04-26T11:00:00Z",
    "title": "Breaking: New Tech Release",
    "description": "A groundbreaking technology was just released...",
    "published_at": "2025-04-26T10:50:00Z",
    "url": "https://example.com/article/123",
    "feed_id": "660e8400-e29b-41d4-a716-446655440000"
  }
]
```

**Status Code:** `200 OK`

**Error Response:**

```json
{
  "error": "Couldn't get posts for user: database error"
}
```

**Status Code:** `400 Bad Request`

---

## 🚀 Installation & Setup

### Prerequisites

- Go 1.25.6 or higher
- Docker and Docker Compose
- PostgreSQL 16 (or use Docker)
- Git

### Step 1: Clone the Repository

```bash
git clone https://github.com/alinaqigit/rss-generator-project.git
cd rss-generator-project
```

### Step 2: Install Go Dependencies

```bash
go mod download
go mod verify
```

### Step 3: Set Up Environment Variables

Create a `.env` file in the project root:

```bash
cp .env.example .env
# Edit .env with your values
```

See [Environment Variables](#-environment-variables) section for details.

### Step 4: Start PostgreSQL Database

```bash
docker-compose up -d postgres
```

This starts a PostgreSQL 16 container with proper health checks.

### Step 5: Run Database Migrations

```bash
# Install goose migration tool
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir ./sql/migrations postgres "$DATABASE_URL" up
```

### Step 6: Generate Database Code (Optional)

If you modify SQL queries:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code from queries
sqlc generate
```

---

## 🔐 Environment Variables

Create a `.env` file with the following variables:

```bash
# Server Configuration
PORT=8080                           # HTTP server port
HOST=127.0.0.1                      # Server host/IP

# Database Configuration
DATABASE_URL="postgres://postgres:postgres@localhost:5544/rss_generator?sslmode=disable"

# Or individually:
# DB_HOST=localhost
# DB_PORT=5544
# DB_USER=postgres
# DB_PASSWORD=postgres
# DB_NAME=rss_generator
```

### Variable Descriptions

| Variable       | Description                  | Example                             | Required                      |
| -------------- | ---------------------------- | ----------------------------------- | ----------------------------- |
| `PORT`         | HTTP server listening port   | `8080`                              | ✅ Yes                        |
| `HOST`         | Server binding address       | `0.0.0.0` or `127.0.0.1`            | ✅ Yes                        |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host:port/db` | ✅ Yes                        |
| `DB_USER`      | PostgreSQL username          | `postgres`                          | ❌ No (if using DATABASE_URL) |
| `DB_PASSWORD`  | PostgreSQL password          | `postgres`                          | ❌ No (if using DATABASE_URL) |
| `DB_HOST`      | PostgreSQL host              | `localhost`                         | ❌ No (if using DATABASE_URL) |
| `DB_PORT`      | PostgreSQL port              | `5432` or `5544`                    | ❌ No (if using DATABASE_URL) |
| `DB_NAME`      | PostgreSQL database name     | `rss_generator`                     | ❌ No (if using DATABASE_URL) |

---

## 🏃 Running the Project

### Option 1: Using Docker Compose (Recommended)

1. **Start all services:**

   ```bash
   docker-compose up -d
   ```

   This starts:
   - PostgreSQL database
   - (Optional) Go application (uncomment in docker-compose.yml)

2. **Check logs:**

   ```bash
   docker-compose logs -f
   ```

3. **Stop services:**
   ```bash
   docker-compose down
   ```

### Option 2: Running Locally

1. **Start the database:**

   ```bash
   docker-compose up -d postgres
   ```

2. **Run migrations:**

   ```bash
   goose -dir ./sql/migrations postgres "$DATABASE_URL" up
   ```

3. **Run the application:**

   ```bash
   go run main.go handler_*.go models.go json.go rss.go scraper.go middleaare_auth.go
   ```

   Or build and run:

   ```bash
   go build -o bin/rss-generator
   ./bin/rss-generator
   ```

4. **Test the API:**

   ```bash
   # Create a user
   curl -X POST http://localhost:8080/v1/user \
     -H "Content-Type: application/json" \
     -d '{"name": "testuser"}'

   # Check health
   curl http://localhost:8080/v1/healthz
   ```

---

## 🛠️ Development

### Project Code Generation

**sqlc** is used for type-safe database code generation from SQL queries.

- **Query Files**: `sql/queries/*.sql`
- **Generated Code**: `internal/db/*.sql.go` (auto-generated, don't edit)
- **Configuration**: `sqlc.yaml`

To regenerate:

```bash
sqlc generate
```

### Database Migrations

Migrations are managed with **Goose**.

- **Location**: `sql/migrations/`
- **Format**: Numbered SQL files with `-- +goose up` and `-- +goose down`
  directives

To create a new migration:

```bash
goose -dir ./sql/migrations postgres "$DATABASE_URL" create migration_name sql
```

### Adding New Endpoints

1. Create a handler function in appropriate `handler_*.go` file
2. Register route in `main.go`
3. If database queries needed, add to `sql/queries/*.sql`
4. Run `sqlc generate` to create query functions

### Code Style

- Follow Go conventions (golangci-lint)
- Use descriptive variable names
- Comment exported functions
- Use error wrapping with context

### Running Tests

```bash
go test ./... -v
```

---

## 📦 Deployment

### Building Docker Image

```dockerfile
# Dockerfile
FROM golang:1.25.6-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o rss-generator .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rss-generator .
EXPOSE 8080
CMD ["./rss-generator"]
```

Build and push:

```bash
docker build -t your-registry/rss-generator:latest .
docker push your-registry/rss-generator:latest
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rss-generator
spec:
  replicas: 2
  selector:
    matchLabels:
      app: rss-generator
  template:
    metadata:
      labels:
        app: rss-generator
    spec:
      containers:
        - name: rss-generator
          image: your-registry/rss-generator:latest
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
            - name: HOST
              value: "0.0.0.0"
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: db-credentials
                  key: connection-string
          livenessProbe:
            httpGet:
              path: /v1/healthz
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /v1/healthz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: rss-generator-service
spec:
  selector:
    app: rss-generator
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: LoadBalancer
```

### Cloud Deployments

#### Azure App Service

```bash
az webapp create --resource-group myGroup --plan myPlan --name rss-gen --runtime "GO|1.25"
```

#### AWS Elastic Beanstalk

```bash
eb init -p go rss-generator
eb create rss-gen-env
eb deploy
```

#### Heroku

```bash
heroku apps:create rss-gen
git push heroku main
```

---

## 📊 Monitoring & Logging

### Health Check

```bash
curl http://localhost:8080/v1/healthz
```

### Application Logs

Check console output for:

- Server startup messages
- Database connection status
- Scraping operations
- Errors and warnings

### Database Connection Issues

```bash
# Test database connection
psql $DATABASE_URL
```

---

## 🐛 Common Issues & Fixes

### Issue: "DATABASE_URL is not set"

**Solution:** Ensure `.env` file exists and `DATABASE_URL` is properly set.

### Issue: "Failed to connect to database"

**Solution:**

- Check PostgreSQL container is running: `docker-compose ps`
- Verify connection string
- Test connection: `psql $DATABASE_URL`

### Issue: "Malformed auth Header"

**Solution:** Ensure API key format is correct: `Authorization: ApiKey {key}`

### Issue: Migrations not applying

**Solution:**

- Verify goose is installed
- Check migration files syntax
- Manually check DB schema: `psql $DATABASE_URL -c "\dt"`

---

## 📈 Performance Optimization

### RSS Scraping

- Adjust concurrency in `startScraping()` (currently 10 goroutines)
- Modify time interval between scrape cycles (currently 1 minute)
- Implement feed priority/rate limiting

### Database Queries

- Add indexes on frequently queried columns
- Use connection pooling
- Monitor slow queries

### API Response Times

- Implement caching
- Add pagination for large result sets
- Use database query optimization

---

## 🔮 Future Enhancements

- [ ] User authentication with JWT tokens
- [ ] Feed categorization and tagging
- [ ] Advanced search and filtering
- [ ] Feed recommendation engine
- [ ] Email notifications for new posts
- [ ] Web UI for feed management
- [ ] Rate limiting and API throttling
- [ ] Full-text search on posts
- [ ] Feed content parsing and enrichment
- [ ] User preferences and settings
- [ ] Analytics and usage statistics

---

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for
details.

---

## 👥 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📧 Support & Contact

- **Issues**: Report via GitHub Issues
- **Discussions**: Use GitHub Discussions for Q&A
- **Email**: alinaqigit@example.com

---

## 🙏 Acknowledgments

- Go community for excellent tooling
- PostgreSQL for reliability
- Chi router for simplicity
- Docker for containerization

---

**Last Updated:** April 26, 2025  
**Version:** 1.0.0
