# Architecture Guide

This document provides an in-depth look at the RSS Generator Project
architecture, design patterns, and system design decisions.

## Table of Contents

- [System Architecture](#system-architecture)
- [Component Overview](#component-overview)
- [Data Flow](#data-flow)
- [Design Patterns](#design-patterns)
- [Security Architecture](#security-architecture)
- [Scalability Considerations](#scalability-considerations)

---

## System Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       External Systems                       │
│  • RSS Feed Publishers (external RSS URLs)                   │
│  • Client Applications (mobile, web, desktop)                │
└────────────────────────────┬────────────────────────────────┘
                             │
        ┌────────────────────┴─────────────────────┐
        │                                          │
        ▼                                          ▼
┌──────────────────┐                   ┌──────────────────────┐
│   HTTP Clients   │                   │  RSS Feed Scrapers   │
│   (Requests)     │                   │    (External URLs)   │
└────────┬─────────┘                   └──────────┬───────────┘
         │                                        │
         ▼                                        ▼
    ┌────────────────────────────────────────────────┐
    │        API Server (Go + Chi Router)            │
    │  ┌──────────────────────────────────────────┐  │
    │  │    HTTP Request Handlers                 │  │
    │  │  • User Management                       │  │
    │  │  • Feed Management                       │  │
    │  │  • Feed Follow Management                │  │
    │  │  • Post Retrieval                        │  │
    │  └──────────────────────────────────────────┘  │
    │  ┌──────────────────────────────────────────┐  │
    │  │    Middleware Stack                      │  │
    │  │  • CORS Handler                          │  │
    │  │  • API Key Authentication                │  │
    │  │  • Error Handling                        │  │
    │  └──────────────────────────────────────────┘  │
    │  ┌──────────────────────────────────────────┐  │
    │  │    Background Services (Goroutines)      │  │
    │  │  • RSS Scraper                           │  │
    │  │    - Concurrent feed fetching            │  │
    │  │    - XML parsing                         │  │
    │  │    - Post persistence                    │  │
    │  └──────────────────────────────────────────┘  │
    └────────────┬───────────────────────────────────┘
                 │
                 │ SQL Queries (sqlc generated)
                 ▼
    ┌────────────────────────────────────────┐
    │   Database Layer (sqlc)                 │
    │  • Query abstraction                    │
    │  • Type-safe prepared statements        │
    │  • Auto-generated database code         │
    └────────────┬───────────────────────────┘
                 │
                 ▼ libpq driver
    ┌────────────────────────────────────────┐
    │   PostgreSQL Database (Port 5544)      │
    │  • Users                                │
    │  • Feeds                                │
    │  • Feed Follows (junction table)        │
    │  • Posts                                │
    └────────────────────────────────────────┘
```

---

## Component Overview

### 1. API Server (Main Application)

**File:** `main.go`

**Responsibilities:**

- Initialize HTTP server with Chi router
- Load environment variables
- Establish database connection
- Register API routes
- Start background RSS scraper
- Handle graceful shutdown

**Key Configuration:**

```go
apiConfig struct {
    DB *db.Queries  // Database query interface
}
```

### 2. Request Handlers

**Files:**

- `handler_user.go` - User CRUD operations
- `handler_feed.go` - Feed management
- `handler_feedFollow.go` - Feed subscription
- `handler-readiness.go` - Health check

**Pattern:**

```go
func (apiCfg *apiConfig) handlerName(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    // 2. Validate input
    // 3. Execute database operation
    // 4. Return response
}
```

**Authentication Handlers:** Some handlers wrap with auth middleware:

```go
func (apiCfg *apiConfig) handlerName(w http.ResponseWriter, r *http.Request, user db.User) {
    // Receives authenticated user from middleware
}
```

### 3. Authentication Middleware

**File:** `middleware_auth.go`

**Flow:**

```
HTTP Request with Authorization Header
        │
        ▼
Extract API Key from Header
        │
        ├─ Missing? → Return 403 Forbidden
        ├─ Invalid format? → Return 403 Forbidden
        │
        ▼
Query Database for User by API Key
        │
        ├─ Not found? → Return 400 Bad Request
        │
        ▼
Attach User Object to Handler
        │
        ▼
Execute Protected Handler
```

**API Key Format:**

```
Authorization: ApiKey {64-character-hex-string}
```

### 4. Database Layer

**Files:**

- `internal/db/db.go` - Database initialization
- `internal/db/models.go` - Database models
- `internal/db/*.sql.go` - Auto-generated query functions

**Technology:** sqlc (SQL Compiler)

**Benefits:**

- Type-safe database queries
- SQL queries written in `.sql` files
- Automatic code generation
- No runtime reflection

### 5. RSS Scraper Background Service

**File:** `scraper.go`

**Architecture:**

```go
startScraping(db, concurrency=10, interval=1min)
    │
    ├─ Start ticker → Fire every interval
    │
    └─ For each tick:
        ├─ Get next N feeds to fetch (where N = concurrency)
        ├─ Create WaitGroup for concurrent processing
        │
        └─ For each feed (in parallel goroutine):
            ├─ Fetch RSS feed from URL (HTTP GET)
            ├─ Parse XML structure
            ├─ For each item in feed:
            │   ├─ Extract title, URL, description, publish date
            │   ├─ Create Post record in database
            │   └─ Handle duplicate URLs gracefully
            ├─ Update feed's last_fetched_at timestamp
            └─ Mark feed as fetched

        Wait for all goroutines to complete
```

**Concurrency Model:**

- Configurable goroutine pool (default: 10)
- WaitGroup ensures all scrapes complete before next tick
- Non-blocking asyncio operation

### 6. Data Models

**File:** `models.go`

**Conversion Pattern:**

```
Database Models (db.User, db.Feed, etc.)
        │
        └─ Conversion Functions
              (database_X_to_X)
        │
        ▼
API Models (User, Feed, etc.)
        │
        ▼
JSON Serialization
        │
        ▼
HTTP Response
```

**Rationale:**

- Separates database concerns from API contract
- Allows API schema evolution independently
- Type-safe transformations

---

## Data Flow

### 1. User Creation Flow

```
POST /v1/user
  ├─ Parse JSON body
  ├─ Validate username
  ├─ DB.CreateUser()
  │   ├─ INSERT INTO users (id, username, api_key)
  │   ├─ Generate UUID
  │   └─ Auto-generate API key (SHA256 hash)
  ├─ Convert db.User → User
  └─ Return JSON response [201 Created]
```

### 2. Feed Subscription Flow

```
POST /v1/feed-follow
  ├─ Middleware: Extract API key, get User
  ├─ Parse JSON body (feed_id)
  ├─ Validate feed exists
  ├─ DB.CreateFeedFollow()
  │   ├─ INSERT INTO feed_follows
  │   ├─ UNIQUE (user_id, feed_id) constraint
  │   └─ Return relationship
  ├─ Convert db.FeedFollow → FeedFollow
  └─ Return JSON response [201 Created]
```

### 3. Post Retrieval Flow

```
GET /v1/posts
  ├─ Middleware: Extract API key, get User
  ├─ DB.GetPostsForUser(user_id, limit=10)
  │   ├─ SELECT * FROM posts WHERE feed_id IN (
  │   │     SELECT feed_id FROM feed_follows
  │   │     WHERE user_id = ?
  │   │   )
  │   ├─ ORDER BY published_at DESC
  │   ├─ LIMIT 10
  │   └─ Return posts slice
  ├─ Convert []db.Post → []Post
  └─ Return JSON response [200 OK]
```

### 4. RSS Scraping Flow

```
Background Goroutine (started at app init)
  │
  ├─ Every 1 minute:
  │   ├─ DB.GetNextFeedsToBeFetched(limit=10)
  │   │   ├─ SELECT * FROM feeds
  │   │   ├─ WHERE last_fetched_at IS NULL
  │   │   │   OR last_fetched_at < NOW() - INTERVAL '...',
  │   │   ├─ ORDER BY last_fetched_at ASC
  │   │   └─ LIMIT 10
  │   │
  │   ├─ For each feed (concurrent goroutines):
  │   │   ├─ urlToFeed(feed.url)
  │   │   │   ├─ HTTP GET feed.url (timeout: 10s)
  │   │   │   ├─ Parse XML to RSSFeed struct
  │   │   │   └─ Return parsed feed
  │   │   │
  │   │   ├─ For each RSSItem in feed:
  │   │   │   ├─ Extract: title, link, description, pubDate
  │   │   │   ├─ DB.CreatePost() with feed_id
  │   │   │   │   ├─ INSERT INTO posts (...)
  │   │   │   │   └─ Handle UNIQUE constraint on url
  │   │   │   └─ Log success/failure
  │   │   │
  │   │   ├─ DB.MarkFeedAsFetched(feed.id)
  │   │   │   ├─ UPDATE feeds
  │   │   │   ├─ SET last_fetched_at = NOW()
  │   │   │   └─ WHERE id = feed.id
  │   │   │
  │   │   └─ Signal WaitGroup completion
  │   │
  │   └─ Wait for all goroutines to finish (WaitGroup.Wait())
  │
  └─ Repeat next cycle after interval
```

---

## Design Patterns

### 1. Receiver Pattern (Dependency Injection)

```go
// apiConfig acts as dependency container
type apiConfig struct {
    DB *db.Queries
}

// Methods use receiver syntax for dependency access
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
    apiCfg.DB.CreateUser(...)
}
```

**Benefit:**

- Cleaner than global variables
- Easy to test with mock DB
- Follows Go idioms

### 2. Type-Safe Database Queries (sqlc)

```sql
-- sql/queries/users.sql
-- name: CreateUser :one
INSERT INTO users (id, username)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, username, api_key;
```

Generates:

```go
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
    row := q.queryRow(ctx, createUser, arg.ID, arg.Username)
    var i User
    err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt, &i.Username, &i.ApiKey)
    return i, err
}
```

**Benefit:**

- Compile-time query validation
- No runtime reflection overhead
- Type safety

### 3. Authentication Wrapper Middleware

```go
type authHandler func(http.ResponseWriter, *http.Request, db.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Extract API key
        // 2. Lookup user
        // 3. Call handler with user context
    }
}
```

**Pattern:**

- `authHandler` receives authenticated user
- Middleware handles authentication logic
- Composition of handlers

### 4. Error Handling

```go
func responseWithError(w http.ResponseWriter, code int, msg string) {
    type errResponse struct {
        Error string `json:"error"`
    }
    responseWithJson(w, code, errResponse{Error: msg})
}
```

**Standardized errors:**

- 400 Bad Request: Client errors, validation failures
- 403 Forbidden: Authentication failures
- 404 Not Found: Resource not found
- 500 Internal Server Error: Server errors
- 204 No Content: Success with no return value

### 5. Concurrent Worker Pool (Scraper)

```go
wg := &sync.WaitGroup{}
for _, feed := range feeds {
    wg.Add(1)
    go scrapFeed(db, wg, feed)  // Non-blocking
}
wg.Wait()  // Block until all complete
```

**Benefits:**

- Bounded concurrency (10 workers)
- Synchronization across goroutines
- Resource-efficient

---

## Security Architecture

### 1. API Key Authentication

**Generation:**

- Random 256-bit value
- SHA256 hash encoded to hex
- Generated automatically on user creation
- Stored in database (unique constraint)

**Validation:**

```
Authorization Header → Extract Key → Query Database → Return User
```

**Storage:**

- Stored as 64-character hex string
- Indexed for fast lookup
- Unique constraint prevents duplicates

### 2. CORS Configuration

```go
cors.Options{
    AllowedOrigins:   []string{"https://*", "http://*"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    AllowCredentials: true,
}
```

**Security Implications:**

- Allows all origins (consider restricting in production)
- Credential support enabled
- Standard HTTP methods allowed

### 3. Database Security

**Features:**

- Foreign keys with CASCADE delete
- UNIQUE constraints (username, api_key, url)
- Parameterized queries (SQL injection prevention via sqlc)
- Timestamps (created_at, updated_at) for audit trail

**Recommendations:**

- Enable PostgreSQL SSL
- Use connection string with sslmode=require
- Implement row-level security (RLS)
- Regular backup strategy

### 4. Input Validation

**Current Level:**

- JSON parsing catches malformed data
- Database constraints enforce uniqueness
- URL validation through HTTP client timeout

**Recommendations:**

- Add explicit validation layer
- Sanitize string inputs
- Validate URLs before storage
- Rate limiting on API endpoints

---

## Scalability Considerations

### 1. Horizontal Scaling

**Current Bottlenecks:**

- Single database connection
- RSS scraper runs on single instance
- No caching layer

**Solutions:**

```
┌─────────────────┐
│  Load Balancer  │
│   (nginx/HAProxy)
└────────┬────────┘
         │
    ┌────┴────┬────────┬─────────┐
    ▼         ▼        ▼         ▼
┌────────┐┌────────┐┌────────┐┌────────┐
│App-1  ││App-2  ││App-3  ││App-4  │
└────┬───┘└────┬───┘└────┬───┘└────┬───┘
     │         │        │        │
     └─────────┼────────┼─────────┘
               ▼
        ┌──────────────┐
        │ PostgreSQL   │
        │ (Connection  │
        │  Pool)       │
        └──────────────┘
```

### 2. Vertical Scaling

**Performance Improvements:**

- Increase RSS scraper concurrency (currently 10)
- Adjust scrape interval (currently 1 minute)
- Increase database connection pool size
- Database query optimization (indexes)

### 3. Caching Strategy

```
Request
  │
  ├─ Check Redis Cache
  │   ├─ HIT → Return cached response
  │   └─ MISS → Continue
  │
  ├─ Query Database
  │
  ├─ Cache Result (TTL: 5 minutes)
  │
  └─ Return Response
```

**Candidates:**

- User profiles (rarely change)
- Feed list (read-heavy)
- Posts (read-only after creation)

### 4. Database Optimization

**Index Opportunities:**

```sql
-- Improve scraper feed fetch query
CREATE INDEX idx_feeds_last_fetched_at ON feeds(last_fetched_at);

-- Improve post retrieval
CREATE INDEX idx_posts_feed_id ON posts(feed_id);
CREATE INDEX idx_posts_published_at ON posts(published_at DESC);

-- Improve feed_follow lookup
CREATE INDEX idx_feed_follows_user_id ON feed_follows(user_id);
CREATE INDEX idx_feed_follows_feed_id ON feed_follows(feed_id);
```

### 5. Message Queue Integration

```
HTTP Request
  │
  ├─ Validate
  ├─ Store metadata
  └─ Enqueue to Message Queue (RabbitMQ/Kafka)
        │
        └─ Return 202 Accepted to client

Background Workers (Multiple)
  ├─ Dequeue task
  ├─ Process (create user, feed, etc.)
  ├─ Update status
  └─ Cleanup
```

---

## Monitoring & Observability

### 1. Logging Strategy

```go
log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)
log.Println("Error fetching feeds: ", err)
log.Println("Error marking a feed as fetched: ", err)
```

**Improvement:**

- Structured logging (JSON format)
- Log levels (DEBUG, INFO, WARN, ERROR)
- Centralized log aggregation

### 2. Metrics to Track

- API response times (histogram)
- Request count by endpoint (counter)
- Error rates by type (counter)
- Database query duration (histogram)
- Active scraper goroutines (gauge)
- Feed scrape success/failure rate (counter)
- Posts created per cycle (counter)

### 3. Health Checks

```
GET /v1/healthz → 200 OK {}
```

**Enhanced health check:**

- Database connectivity
- Redis connectivity
- Active goroutine count
- Memory usage
- Uptime

---

## Deployment Architecture

### Development

```
┌─────────────────────────────┐
│    Local Development        │
├─────────────────────────────┤
│ localhost:8080              │
│ postgres://localhost:5544   │
│ Hot reload (optional)       │
└─────────────────────────────┘
```

### Production

```
┌────────────────────────────────────────────────┐
│          CDN / DDoS Protection                 │
│          (Cloudflare, AWS Shield)              │
└─────────────────────┬──────────────────────────┘
                      │
        ┌─────────────┴──────────────┐
        ▼                            ▼
┌──────────────────┐      ┌──────────────────┐
│  API Server 1    │      │  API Server 2    │
│  (ECS / K8s)     │      │  (ECS / K8s)     │
└────────┬─────────┘      └────┬─────────────┘
         │                     │
         └─────────────┬───────┘
                       ▼
        ┌─────────────────────────┐
        │  Database (RDS/CloudSQL)│
        │  - Read Replicas        │
        │  - Backups              │
        │  - Point-in-time restore│
        └─────────────────────────┘
```

---

## Future Architectural Enhancements

1. **Event-Driven Architecture:** Replace polling with event streaming
2. **Microservices:** Separate scrap service from API service
3. **GraphQL API:** In addition to REST
4. **Real-time Updates:** WebSocket support for live feed updates
5. **ML Integration:** Feed recommendations, duplicate detection
6. **Distributed Tracing:** OpenTelemetry integration
7. **Feature Flags:** A/B testing and gradual rollouts

---

**Last Updated:** April 26, 2025
